package connection_test

import (
	"context"
	"fmt"
	"testing"

	nmcli_go "github.com/KunMengcode/nmcli-go"
	"github.com/KunMengcode/nmcli-go/connection"
)

func TestManager_Up(t *testing.T) {
	m := nmcli_go.NewNMCli()
	out, err := m.Connection.Up(context.Background(), "hotspot", connection.UpOptions{})
	if err != nil {
		return
	}
	t.Log(out)
}

func TestManager_Show(t *testing.T) {
	m := nmcli_go.NewNMCli()
	out, err := m.Connection.Show(context.Background(), "hotspot")
	if err != nil {
		return
	}
	t.Log(out)
}

func TestManager_GetInterfaceConnectionType(t *testing.T) {
	m := nmcli_go.NewNMCli()
	out, err := m.Connection.GetConnectionType(context.Background(), "wlp0s20f3")
	if err != nil {
		fmt.Println(out, err)
		return
	}
	t.Log(out)
}
