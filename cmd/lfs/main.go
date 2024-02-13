package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var store = map[string]string{}

func handleError(e error) {
	fmt.Println(e)
	os.Exit(1)
}

func main() {
	file, err := os.OpenFile(".db", os.O_CREATE|os.O_WRONLY, 0744)

	if err != nil {
		handleError(err)
	}
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		handleError(err)
	}
	fmt.Println("Server listening on 8080")
	conn, err := listener.Accept()
	if err != nil {
		handleError(err)
	}

	defer conn.Close()

	for {
		buffer := make([]byte, 100)
		conn.Read(buffer)
		tokens := strings.Fields(string(buffer))
		key := tokens[0]
		val := tokens[1]
		if key == "quit" {
			break
		}
		record := fmt.Sprintf("%s,%s\n", key, val)
		file.Write([]byte(record))
	}
}
