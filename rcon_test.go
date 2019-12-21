package mc_rcon

import "testing"

var conn *MCConn

func TestMCConn_Open(t *testing.T) {
	conn = new(MCConn)
	err := conn.Open("minecraft:25575", "testpw")
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
