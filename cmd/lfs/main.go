package main

import "lfs/internal"

func main() {
	server, err := internal.CreateServer("localhost", 8080, "tcp")
	if err != nil {
		internal.HandleError("Could not start server", err)
	}
	err = server.AcceptConnections(internal.Engine)
	if err != nil {
		internal.HandleError("Could not accept connections", err)
	}
}
