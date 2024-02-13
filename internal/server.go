package internal

import (
	"fmt"
	"net"
)

type server struct {
	address  string
	listener net.Listener
}

func CreateServer(host string, port uint16, networkType string) (*server, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	listener, err := net.Listen(networkType, address)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s server listening on %s", networkType, address)
	serverInstance := server{address, listener}
	return &serverInstance, nil
}

func (s *server) AcceptConnections(handler func(connection net.Conn) error) error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		go handler(conn)
	}
}
