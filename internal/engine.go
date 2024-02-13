package internal

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var index = map[string]int64{}

func Engine(conn *net.Conn, server *server) error {
	var offset int64 = 0
	for {
		buffer := make([]byte, 128)
		n, err := (*conn).Read(buffer)
		if err != nil {
			(*conn).Close()
			HandleError("Could not read from connection", err)
		}
		client := (*conn).RemoteAddr().String()
		server.connections[client] += 1
		file, err := os.OpenFile(".db", os.O_RDWR, 0744)
		if err != nil {
			(*conn).Close()
			HandleError("Could not open file", err)
		}
		// DEBT n-1 because newline is appended when message is sent from netcat on enter.
		msg := string(buffer[:n-1])
		args := strings.Fields(msg)
		cmd := strings.ToLower(args[0])
		switch cmd {
		case "set":
			key := args[1]
			val := args[2]
			handleSet(key, val, file, &offset)
		}
	}
}

func handleSet(key string, val string, file *os.File, offset *int64) {
	defer file.Close()
	record := fmt.Sprintf("%s,%s\n", key, val)
	serializedRecord := []byte(record)
	_, err := file.Seek(*offset, 0)
	if err != nil {
		HandleError("Could not seek to file", err)
	}
	bytes, err := file.Write(serializedRecord)
	index[key] = *offset
	*offset = *offset + int64(bytes)
	if err != nil {
		HandleError("Could not write to DB", err)
	}
}
