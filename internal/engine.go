package internal

import (
	"bufio"
	"errors"
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
			if err.Error() != "EOF" {
				fmt.Println("Could not read from connection\n", err)
			}
			return nil
		}
		client := (*conn).RemoteAddr().String()
		server.connections[client] += 1
		// DEBT n-1 because newline is appended when message is sent from netcat on enter.
		msg := string(buffer[:n-1])
		args := strings.Fields(msg)
		cmd := strings.ToLower(args[0])
		switch cmd {
		case "set":
			key := args[1]
			val := args[2]
			response, err := handleSet(key, val)
			if err != nil {
				fmt.Println(err)
				continue
			}
			(*conn).Write([]byte(response + "\n"))
		case "get":
			key := args[1]
			response, err := handleGet(key)
			if err != nil {
				fmt.Println(err)
				continue
			}
			(*conn).Write([]byte(response + "\n"))
		default:
			(*conn).Write([]byte("Incorrect command\n"))
		}
	}
}

func handleGet(key string) (string, error) {
	file, err := os.OpenFile(".db", os.O_RDONLY, 0744)
	if err != nil {
		return "", errors.New("Could not open file for reading\n" + err.Error())
	}
	defer file.Close()
	offset := index[key]
	file.Seek(offset, 0)
	reader := bufio.NewReader(file)
	buffer, err := reader.ReadBytes('\n')
	if err != nil {
		return "", errors.New("Could not read from DB\n" + err.Error())
	}
	record := string(buffer[:len(buffer)-1])
	pair := strings.Split(record, ",")
	return pair[1], nil
}

func handleSet(key string, val string) (string, error) {
	file, err := os.OpenFile(".db", os.O_WRONLY, 0744)
	if err != nil {
		return "", errors.New("Could not open file for writing\n" + err.Error())
	}
	defer file.Close()
	record := fmt.Sprintf("%s,%s\n", key, val)
	serializedRecord := []byte(record)
	offset, err := file.Seek(0, 2)
	if err != nil {
		return "", errors.New("Could not seek to file\n" + err.Error())
	}
	_, err = file.Write(serializedRecord)
	index[key] = offset
	if err != nil {
		return "", errors.New("Could not write to DB\n" + err.Error())
	}
	return "OK", nil
}
