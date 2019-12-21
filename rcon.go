package mc_rcon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

// Copied pretty much from github.com/bearbin/mcgorcon, but with some minor modifications

const (
	BadAuth        = -1
	PayloadMaxSize = 1460
)

const (
	PacketResponse = iota
	_
	PacketCommand
	PacketLogin
)

type MCConn struct {
	conn     net.Conn
	password string
}

type packetType int32

type RCONHeader struct {
	Size      int32
	RequestID int32
	Type      packetType
}

func (c *MCConn) Open(addr, password string) error {
	log.Println("Dial " + addr)
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return err
	}

	c.conn = conn
	c.password = password

	return nil
}

func (c *MCConn) Close() {
	_ = c.conn.Close()
}

// SendCommand sends a command to the server and returns the result (often nothing).
func (c *MCConn) SendCommand(command string) (string, error) {
	// Send the packet.
	if len([]byte(command)) > PayloadMaxSize {
		return "", errors.New("payload too large")
	}
	head, payload, err := c.sendPacket(PacketCommand, []byte(command))
	if err != nil {
		return "", err
	}

	// Auth was bad, throw error.
	if head.RequestID == BadAuth {
		return "", errors.New("bad auth, could not send command")
	}
	return string(payload), nil
}

// authenticate authenticates the user with the server.
func (c *MCConn) Authenticate() error {
	// Send the packet.
	head, _, err := c.sendPacket(PacketLogin, []byte(c.password))
	if err != nil {
		return err
	}

	// If the credentials were bad, throw error.
	if head.RequestID == BadAuth {
		return errors.New("bad auth, could not authenticate")
	}

	return nil
}

// sendPacket sends the binary packet representation to the server and returns the response.
func (c *MCConn) sendPacket(t packetType, p []byte) (RCONHeader, []byte, error) {
	// Generate the binary packet.
	packet, err := packetise(t, p)
	if err != nil {
		return RCONHeader{}, nil, err
	}

	// Send the packet over the wire.
	_, err = c.conn.Write(packet)
	if err != nil {
		return RCONHeader{}, nil, err
	}
	// Receive and decode the response.
	head, payload, err := depacketise(c.conn)
	if err != nil {
		return RCONHeader{}, nil, err
	}

	return head, payload, nil
}

// packetise encodes the packet type and payload into a binary representation to send over the wire.
func packetise(t packetType, p []byte) ([]byte, error) {
	// Generate a random request ID.
	pad := [2]byte{}
	length := int32(len(p) + 10)
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, length)
	_ = binary.Write(&buf, binary.LittleEndian, int32(0))
	_ = binary.Write(&buf, binary.LittleEndian, t)
	_ = binary.Write(&buf, binary.LittleEndian, p)
	_ = binary.Write(&buf, binary.LittleEndian, pad)
	// Notchian server doesn't like big packets :(
	if buf.Len() >= 1460 {
		return nil, errors.New("packet too big when packetising")
	}
	// Return the bytes.
	return buf.Bytes(), nil
}

// depacketise decodes the binary packet into a native Go struct.
func depacketise(r io.Reader) (RCONHeader, []byte, error) {
	head := RCONHeader{}
	err := binary.Read(r, binary.LittleEndian, &head)
	if err != nil {
		return RCONHeader{}, nil, err
	}
	payload := make([]byte, head.Size-8)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return RCONHeader{}, nil, err
	}

	// Some basic sanity checking
	if head.Type != PacketResponse && head.Type != PacketCommand {
		return RCONHeader{}, nil, errors.New("bad packet type")
	}
	return head, payload[:len(payload)-2], nil
}
