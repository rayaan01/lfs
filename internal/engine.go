package internal

import (
	"fmt"
	"net"
)

func Engine(conn *net.Conn, server *server) error {
	for {
		buffer := make([]byte, 1024)
		n, err := (*conn).Read(buffer)
		if err != nil {
			(*conn).Close()
			return err
		}
		// n-1 because newline is appended when message is sent from netcat
		msg := string(buffer[:n-1])
		remoteAddr := (*conn).RemoteAddr().String()
		server.connections[remoteAddr] += 1
		fmt.Println(msg)
	}
}
