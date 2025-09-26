package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ZanyLeonic/mineprimer/protocol"
)

const (
	ListenPort = ":25565"
)

func main() {
	s, err := net.Listen("tcp4", ListenPort)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	for {
		conn, err := s.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go HandleConnection(conn)
	}
}

func HandleConnection(c net.Conn) {
	defer c.Close()
	
	c.SetDeadline(time.Now().Add(10 * time.Second))
	log.Printf("New connection from %s\n", c.RemoteAddr())

	pkt, err := protocol.ReadPacket(c)
	if err != nil {
		return
	}

	if pkt.ID != 0x00 {
		log.Printf("First Packet not handshake (id=0x%02X)", pkt.ID)
		return
	}

	next, err := HandleHandshake(pkt)
	if err != nil {
		log.Printf("Handshake parse error: %v", err)
		return
	}

	log.Printf("State received -> %v\n", next)

	switch next.NextState {
	case protocol.StateStatus:
		HandleStatusState(c, next)
	}
}

func HandleHandshake(pkt *protocol.Packet) (protocol.HandshakeInfo, error) {
	buf := bytes.NewReader(pkt.Payload)
	handshakeInfo := protocol.HandshakeInfo{}

	protoVer, err := protocol.ReadVarInt(buf)
	if err != nil {
		return handshakeInfo, err
	}

	addr, err := protocol.ReadString(buf)
	if err != nil {
		return handshakeInfo, err
	}

	port, err := protocol.ReadUnsignedShort(buf)
	if err != nil {
		return handshakeInfo, err
	}

	nextState, err := protocol.ReadVarInt(buf)
	if err != nil {
		return handshakeInfo, err
	}

	handshakeInfo = protocol.HandshakeInfo{
		ProtocolVersion: protoVer,
		Address: addr,
		Port: port,
		NextState: protocol.ConnectionState(nextState),
	}

	log.Printf("Handshake proto=%d addr=%s port=%d next=%d\n", protoVer, addr, port, nextState)
	if nextState == 1 {
		return handshakeInfo, nil
	} else if nextState == 2 {
		return handshakeInfo, nil
	} 

	return handshakeInfo, fmt.Errorf("unknown next state %d", nextState)
}

func HandleStatusState(c net.Conn, h protocol.HandshakeInfo) {
	for {
		pkt, err := protocol.ReadPacket(c)
		if err != nil {
			log.Printf("Cannot read status: %v\n", err)
			return
		}

		switch pkt.ID {
		case 0x00:
			status := protocol.PingStatus{
				Version: protocol.StatusVersion{
					Name: "Server on Standby", Protocol: int(h.ProtocolVersion),
				}, 
				Players: protocol.StatusPlayerInfo{
					Max: 20, 
					Online: 0,
				}, 
				Description: protocol.StatusDescription{Text: "Server is on standby"},
			}

			b, err := json.Marshal(status)
			if err != nil {
				log.Printf("Cannot marshal Status JSON: %v\n", err)
				return
			}

			var payload bytes.Buffer

			protocol.WriteString(&payload, protocol.String(b))
			protocol.WritePacket(c, 0x00, payload.Bytes())

			log.Printf("Sent response %s\n", string(b))
		case 0x01:
			protocol.WritePacket(c, 0x01, pkt.Payload)
			log.Printf("ping len=%d\n", len(pkt.Payload))
		default:
			log.Printf("Unknown status packet id=0x%02X", pkt.ID)
		}
	}
}