package connection

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/KunMengcode/nmcli-go/utils"
)

type UpOptions struct {
	Ifname      string `json:"ifname"`
	BSSID       string `json:"ap"`
	Passwd_File string `json:"passwd-file"`
}

func (m Manager) Up(ctx context.Context, ID string, args UpOptions) (string, error) {
	cmdArgs := []string{"connection", "up", ID}
	cmdArgs = append(cmdArgs, utils.Marshal(args)...)
	output, err := m.CommandContext(ctx, nmcliCmd, cmdArgs...).Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute nmcli with args %+q: %w", cmdArgs, err)
	}
	return string(output), nil
}

func (m Manager) Modify(ctx context.Context, temporary bool, ID string, option map[string]string) (string, error) {
	cmdArgs := []string{"connection", "modify"}
	if temporary {
		cmdArgs = append(cmdArgs, "--temporary")
	}
	cmdArgs = append(cmdArgs, ID)
	for k, v := range option {
		cmdArgs = append(cmdArgs, k, v)
	}
	output, err := m.CommandContext(ctx, nmcliCmd, cmdArgs...).Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute nmcli with args %+q: %w", cmdArgs, err)
	}
	return string(output), nil
}

func (m Manager) Show(ctx context.Context, ConnId string) (map[string][][]string, error) {
	cmdArgs := []string{"-s", "-g", "all", "connection", "show", ConnId}
	output, err := m.CommandContext(ctx, nmcliCmd, cmdArgs...).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute nmcli with args %+q: %w", cmdArgs, err)
	}
	return utils.ParseCmdHaveFieldNameOutput(output), nil
}

// GetConnectionType 获取指定网卡接口的连接类型
func (m Manager) GetConnectionType(ctx context.Context, network string) (string, error) {
	deviceShowCmd := exec.CommandContext(ctx, "nmcli", "device", "show", network)
	grepCmd := exec.CommandContext(ctx, "grep", "GENERAL.CONNECTION")

	grepCmd.Stdin, _ = deviceShowCmd.StdoutPipe()
	var outBuf strings.Builder
	grepCmd.Stdout = &outBuf

	if err := grepCmd.Start(); err != nil {
		return "unknown", fmt.Errorf("failed to start grep command: %w", err)
	}
	if err := deviceShowCmd.Run(); err != nil {
		return "unknown", fmt.Errorf("failed to execute nmcli device show for interface '%s': %w", network, err)
	}

	if err := grepCmd.Wait(); err != nil {
		return "unknown", fmt.Errorf("grep command failed: %w", err)
	}

	outputLine := strings.TrimSpace(outBuf.String())
	parts := strings.Split(outputLine, ":")
	if len(parts) < 2 {
		return "unknown", fmt.Errorf("unexpected output format from nmcli device show: %s", outputLine)
	}
	connName := strings.TrimSpace(parts[1]) // "hotspot"

	if strings.ToLower(connName) == "hotspot" {
		return "hotspot", nil
	} else {
		return "wifi", nil
	}
}
