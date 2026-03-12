package nmcli_go

import (
	"context"
	"os/exec"

	"github.com/KunMengcode/nmcli-go/connection"

	"github.com/KunMengcode/nmcli-go/device"

	"github.com/KunMengcode/nmcli-go/general"
	"github.com/KunMengcode/nmcli-go/utils"
)

type General interface {
	Hostname(ctx context.Context, args general.HostnameArgs) (string, error)
	Permissions(ctx context.Context) ([]general.Permission, error)
}

type Device interface {
	WiFiList(ctx context.Context, args device.WiFiListOptions) ([]device.WiFi, error)
	WiFiConnect(ctx context.Context, BSSID string, args device.WiFiConnectOptions) (string, error)
	WiFiHotspotCreate(ctx context.Context, args device.WiFiHotspotCreateOptions) (string, error)
	Status(ctx context.Context) ([]device.Status, error)
	Show(ctx context.Context, DeviceInterfaceName string) ([]map[string][][]string, error)
}
type Connection interface {
	Up(ctx context.Context, ID string, args connection.UpOptions) (string, error)
	Show(ctx context.Context, ConnId string) (map[string][][]string, error)
	Modify(ctx context.Context, temporary bool, ID string, option map[string]string) (string, error)
	GetConnectionType(ctx context.Context, network string) (string, error)
}

type NMCli struct {
	// should be used to exec custom nmcli commands
	CommandContext func(ctx context.Context, name string, args ...string) utils.Cmd
	General        General
	Device         Device
	Connection     Connection
}

type Option = func(cli *NMCli)

func NewNMCli(opts ...Option) NMCli {
	cli := NMCli{
		CommandContext: func(ctx context.Context, name string, args ...string) utils.Cmd {
			return exec.CommandContext(ctx, name, args...)
		},
	}
	for i := range opts {
		opts[i](&cli)
	}

	cli.General = general.Manager{CommandContext: cli.CommandContext}
	cli.Device = device.Manager{CommandContext: cli.CommandContext}
	cli.Connection = connection.Manager{CommandContext: cli.CommandContext}

	return cli
}
