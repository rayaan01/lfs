package main

import (
	"fmt"
	"lfs/internal"
)

func main() {
	server, err := internal.CreateServer("localhost", 8080, "tcp")
	if err != nil {
		fmt.Println("Could not start server", err)
	}
	server.AcceptConnections(internal.Engine)
}
