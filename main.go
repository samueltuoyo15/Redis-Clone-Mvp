package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	commands "github.com/samueltuoyo15/Redis-Clone-Mvp/commands"
	utils "github.com/samueltuoyo15/Redis-Clone-Mvp/utils"
)

// To read one complete RESP message (array form) and returns slice of argument

func parseRESP(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimRight(line, "\r\n")
	if len(line) == 0 {
		return nil, errors.New("Empty line")
	}

	if line[0] == '*' {
		// Array of bulk strings
		count, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}
		args := make([]string, 0, count)
		for i := 0; i < count; i++ {
			sizeLine, err := reader.ReadString('\n')
			if err != nil {
				return nil, err
			}

			sizeLine = strings.TrimRight(sizeLine, "\r\n")
			if sizeLine == "" || sizeLine[0] != '$' {
				return nil, fmt.Errorf("Expected bulk string, got: %s", sizeLine)
			}

			size, err := strconv.Atoi(sizeLine[1:])
			if err != nil {
				return nil, err
			}

			buf := make([]byte, size+2)
			_, err = io.ReadFull(reader, buf)
			if err != nil {
				return nil, err
			}

			arg := string(buf[:size])
			args = append(args, arg)
		}
		return args, nil
	}

	// Support for inline commands
	parts := strings.Split(line, " ")
	return parts, nil
}

func handleConnection(conn net.Conn, store *commands.Store) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		_ = conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		args, err := parseRESP(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			_, _ = writer.Write(utils.EncodeError("ERR parse error: " + err.Error()))
			_ = writer.Flush()
			continue
		}

		resp := commands.HandleCommand(store, args)
		_, _ = writer.Write(resp)
		_ = writer.Flush()

		if len(args) > 0 && strings.ToUpper(args[0]) == "QUIT" {
			return
		}
	}
}

func main() {
	fmt.Println("Go Redis server is listening on port 6378")

	listener, err := net.Listen("tcp", ":6378")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	store := &commands.Store{
		Data:  make(map[string]string),
		TTL:   make(map[string]time.Time),
		Clean: make(chan struct{}),
	}

	defer store.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("New client connected:", conn.RemoteAddr())
		go handleConnection(conn, store)
	}
}
