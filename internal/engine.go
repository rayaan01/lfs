package internal

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var index = map[string]int64{}

func Engine(conn net.Conn, server *server) error {
	for {
		clientAddress := conn.RemoteAddr().String()
		buffer := make([]byte, 0, 4096)
		bytesRead := 0
		err := read(&buffer, &bytesRead, conn)
		if err != nil {
			if err == io.EOF {
				conn.Close()
				return nil
			}
			fmt.Printf("Could not read from connection %s : %s \n", clientAddress, err)
			conn.Write([]byte("Something went wrong!"))
			continue
		}

		server.connections[clientAddress] += 1
		input := string(buffer[:bytesRead-1])
		args := strings.Fields(input)
		response, err := router(args)
		if err != nil {
			if err == io.EOF {
				conn.Close()
				return nil
			}
			fmt.Printf("Could not handle args %s : %s \n", clientAddress, err)
			continue
		}
		conn.Write(response)
	}
}

func read(buffer *[]byte, bytesRead *int, conn net.Conn) error {
	for {
		chunk := make([]byte, 128)
		n, err := conn.Read(chunk)
		if err != nil {
			return err
		}
		*buffer = append(*buffer, chunk[:n]...)
		*bytesRead += n
		if n < cap(chunk) {
			return nil
		}
	}
}

func router(args []string) ([]byte, error) {
	cmd := strings.ToLower(args[0])
	switch cmd {
	case "set":
		key := args[1]
		val := args[2]
		response, err := handleSet(key, val)
		if err != nil {
			return nil, err
		}
		return append(response, []byte("\n")...), nil
	case "get":
		key := args[1]
		response, err := handleGet(key)
		if err != nil {
			return nil, err
		}
		return append(response, []byte("\n")...), nil
	case "exit":
		return nil, io.EOF
	default:
		return []byte("Available commands: set [key] [value], get [key], exit \n"), nil
	}
}

func handleGet(key string) ([]byte, error) {
	file, err := os.OpenFile(".db", os.O_RDONLY, 0644)
	if err != nil {
		return nil, errors.New("Could not open file for reading\n" + err.Error())
	}
	defer file.Close()
	offset := index[key]
	file.Seek(offset, 0)
	reader := bufio.NewReader(file)
	buffer, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, errors.New("Could not read from DB\n" + err.Error())
	}
	record := string(buffer[:len(buffer)-1])
	pair := strings.Split(record, ",")
	return []byte(pair[1]), nil
}

func handleSet(key string, val string) ([]byte, error) {
	file, err := os.OpenFile(".db", os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.New("Could not open file for writing\n" + err.Error())
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	keyLength := uint16(len(key))
	err = binary.Write(writer, binary.LittleEndian, keyLength)
	if err != nil {
		return nil, errors.New("Could not write key length\n" + err.Error())
	}
	_, err = writer.Write([]byte(key))
	if err != nil {
		return nil, errors.New("Could not write key\n" + err.Error())
	}

	valueLength := uint16(len(val))
	err = binary.Write(writer, binary.LittleEndian, valueLength)
	if err != nil {
		return nil, errors.New("Could not write value length\n" + err.Error())
	}
	_, err = writer.Write([]byte(val))
	if err != nil {
		return nil, errors.New("Could not write value\n" + err.Error())
	}

	offset, err := file.Seek(0, 2)
	if err != nil {
		return nil, errors.New("Could not seek to file\n" + err.Error())
	}
	err = writer.Flush()
	if err != nil {
		return nil, errors.New("Could not write to DB\n" + err.Error())
	}
	index[key] = offset
	return []byte("OK"), nil
}
