package main

import (
	"fmt"
	"lfs/internal"
	"os"
)

func main() {
	server, err := internal.CreateServer("localhost", 8080, "tcp")
	if err != nil {
		fmt.Printf("Could not start server: %s \n", err)
		os.Exit(1)
	}
	server.AcceptConnections(internal.Engine)
}
