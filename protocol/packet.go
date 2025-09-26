package protocol

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

const MaxPacketSize = 3 * 3 * 1024

// ConnectionState for all the states
type ConnectionState int
const (
	StateHandshake ConnectionState = iota
	StateStatus
	StateLogin
	StatePlay
)

type Packet struct {
	ID      int32
	Payload []byte
}

func ReadPacket(conn net.Conn) (*Packet, error) {
	length, err := ReadVarInt(conn)
	if err != nil {
		return nil, err
	}

	if length < 0 {
		return nil, fmt.Errorf("packet with negative length %d", length)
	}

	if length > MaxPacketSize {
		return nil, fmt.Errorf("packet length %d is too large", length)
	}

	data := make([]byte, length)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(data)
	id, err := ReadVarInt(buf)
	if err != nil {
		return nil, err
	}

	remaining, _ := io.ReadAll(buf)
	return &Packet{ID: int32(id), Payload: remaining}, nil
}

func WritePacket(conn net.Conn, id VarInt, payload []byte) error {
	var inner bytes.Buffer
	if err := WriteVarInt(&inner, id); err != nil {
		return err
	}

	if len(payload) > 0 {
		if _, err := inner.Write(payload); err != nil {
			return err
		}
	}

	var outer bytes.Buffer
	if err := WriteVarInt(&outer, VarInt(inner.Len())); err != nil {
		return err
	}

	if _, err := outer.Write(inner.Bytes()); err != nil {
		return err
	}

	_, err := conn.Write(outer.Bytes())
	return err
}