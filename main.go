package main

import (
	"log"
	"net"
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
	log.Println("TODO: implement connection handling")
}