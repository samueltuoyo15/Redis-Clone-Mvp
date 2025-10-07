package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"
)

type Store struct {
	mu    sync.RWMutex
	data  map[string]string
	ttl   map[string]time.Time
	clean chan struct{}
}

func NewStore() *Store {
	s := &Store{
		data:  make(map[string]string),
		ttl:   make(map[string]time.Time),
		clean: make(chan struct{}),
	}
	go s.StartTTLWorker()
	return s
}

func (s *Store) StartTTLWorker() {
	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			s.mu.Lock()
			for k, t := range s.ttl {
				if t.Before(now) {
					delete(s.data, k)
					delete(s.ttl, k)
				}
			}
			s.mu.Unlock()
		case <-s.clean:
			ticker.Stop()
			return
		}
	}
}

func (s *Store) Close() {
	close(s.clean)
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if t, ok := s.ttl[key]; ok && time.Now().After(t) {
		s.mu.RUnlock()
		s.mu.Lock()
		delete(s.data, key)
		delete(s.ttl, key)
		s.mu.Unlock()
		return "", false
	}
	v, ok := s.data[key]
	return v, ok
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			return
		}

		fmt.Printf("Received from %s : %s, conn.RemoteAddr() ", message, conn.RemoteAddr())
		conn.Write([]byte("+PONG\r\n"))
	}
}

func main() {
	fmt.Println("Starting Go Redis server on port 6378")

	listener, err := net.Listen("tcp", ":6378")
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("New client connected:", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
