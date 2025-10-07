package commands

import (
	"strconv"
	"strings"
	"sync"
	"time"

	utils "github.com/samueltuoyo15/Redis-Clone-Mvp/utils"
)

type Store struct {
	mu    sync.RWMutex
	Data  map[string]string
	TTL   map[string]time.Time
	Clean chan struct{}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if t, ok := s.TTL[key]; ok && time.Now().After(t) {
		s.mu.RUnlock()
		s.mu.Lock()
		delete(s.Data, key)
		delete(s.TTL, key)
		s.mu.Unlock()
		return "", false
	}

	v, ok := s.Data[key]
	return v, ok
}

func (s *Store) Set(key, value string, expirySeconds int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Data[key] = value
	if expirySeconds > 0 {
		s.TTL[key] = time.Now().Add(time.Duration(expirySeconds) * time.Second)
	} else {
		delete(s.TTL, key)
	}
}

func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.Data[key]
	if ok {
		delete(s.Data, key)
		delete(s.TTL, key)
	}
	return ok
}

func (s *Store) Close() {
	close(s.Clean)
}

func HandleCommand(store *Store, args []string) []byte {
	if len(args) == 0 {
		return utils.EncodeError("ERR no command given")
	}

	cmd := strings.ToUpper(args[0])

	switch cmd {
	case "PING":
		if len(args) == 2 {
			return utils.EncodeBulkString(args[1])
		}
		return utils.EncodeSimpleString("PONG")

	case "ECHO":
		if len(args) < 2 {
			return utils.EncodeError("ERR wrong number of arguments for 'ECHO' command")
		}
		return utils.EncodeBulkString(args[1])

	case "SET":
		if len(args) < 3 {
			return utils.EncodeError("ERR wrong number of arguments for 'SET' command")
		}
		key := args[1]
		value := args[2]
		expiry := 0

		if len(args) >= 5 && strings.ToUpper(args[3]) == "EX" {
			if sec, err := strconv.Atoi(args[4]); err == nil {
				expiry = sec
			}
		}

		store.Set(key, value, expiry)
		return utils.EncodeSimpleString("OK")

	case "GET":
		if len(args) != 2 {
			return utils.EncodeError("ERR wrong number of arguments for 'GET' command")
		}
		key := args[1]
		val, ok := store.Get(key)
		if !ok {
			return []byte("$-1\r\n")
		}
		return utils.EncodeBulkString(val)

	case "DEL":
		if len(args) < 2 {
			return utils.EncodeError("ERR wrong number of arguments for 'DEL' command")
		}
		count := 0
		for _, key := range args[1:] {
			if store.Del(key) {
				count++
			}
		}
		return utils.EncodeInteger(count)

	case "QUIT":
		return utils.EncodeSimpleString("OK")

	default:
		return utils.EncodeError("ERR unknown command '" + args[0] + "'")
	}
}
