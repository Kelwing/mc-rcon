package mc_rcon

import (
	"fmt"
	"os"
	"testing"
)

var conn *MCConn

func TestMCConn_Open(t *testing.T) {
	conn = new(MCConn)
	addr := fmt.Sprintf("%s:%s", os.Getenv("MINECRAFT_HOST"), os.Getenv("MINECRAFT_PORT"))
	err := conn.Open(addr, "testpw")
	if err != nil {
		t.Error("Open failed", err)
		return
	}
}

func TestMCConn_Authenticate(t *testing.T) {
	err := conn.Authenticate()
	if err != nil {
		t.Error("Auth failed", err)
		return
	}
}

func TestMCConn_SendCommand(t *testing.T) {
	resp, err := conn.SendCommand("seed")
	if err != nil {
		t.Error("Command failed", err)
		return
	}

	t.Log(resp)
}

func TestMCConn_Close(t *testing.T) {
	conn.Close()
}
