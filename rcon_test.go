package mc_rcon

import "testing"

func TestOpen(t *testing.T) {
	conn := new(MCConn)
	err := conn.Open("localhost:25575", "testpw")
	if err != nil {
		t.Error("Open failed", err)
		return
	}

	err = conn.Authenticate()
	if err != nil {
		t.Error("Auth failed", err)
		return
	}

	resp, err := conn.SendCommand("tps")
	if err != nil {
		t.Error("Command failed", err)
	}

	t.Log(resp)
}
