package internal

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var index = map[string]int64{}

func Engine(conn *net.Conn, server *server) error {
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
			response := handleSet(key, val, file)
			(*conn).Write([]byte(response + "\n"))
		case "get":
			key := args[1]
			response := handleGet(key, file)
			(*conn).Write([]byte(response + "\n"))
		}
	}
}

func handleGet(key string, file *os.File) string {
	defer file.Close()
	offset := index[key]
	file.Seek(offset, 0)
	reader := bufio.NewReader(file)
	buffer, err := reader.ReadBytes('\n')
	if err != nil {
		HandleError("Could not read from DB", err)
	}
	record := string(buffer[:len(buffer)-1])
	pair := strings.Split(record, ",")
	return pair[1]
}

func handleSet(key string, val string, file *os.File) string {
	defer file.Close()
	record := fmt.Sprintf("%s,%s\n", key, val)
	serializedRecord := []byte(record)
	offset, err := file.Seek(0, 2)
	if err != nil {
		HandleError("Could not seek to file", err)
	}
	_, err = file.Write(serializedRecord)
	index[key] = offset
	if err != nil {
		HandleError("Could not write to DB", err)
	}
	return "OK"
}
