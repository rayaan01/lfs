package internal

import (
	"fmt"
	"net"
)

type server struct {
	address     string
	listener    net.Listener
	connections map[string]uint16
}

func CreateServer(host string, port uint16, networkType string) (*server, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen(networkType, address)
	connections := map[string]uint16{}
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s server listening on %s \n", networkType, address)
	serverInstance := server{address, listener, connections}
	return &serverInstance, nil
}

func (s *server) AcceptConnections(handler func(conn *net.Conn, server *server) error) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			conn.Close()
			HandleError("Could not accept connection", err)
			continue
		}
		s.connections[conn.RemoteAddr().String()] = 0
		go handler(&conn, s)
	}
}
