package parrotserver

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"sync"

	"go.uber.org/zap"
)

/* Global vars */
var redisServer server /* Server global state */

func GetServerInstance() server {
	return redisServer
}

type server struct {
	logger   *zap.Logger
	errChan  chan error
	listener net.Listener
	dict     map[string]string
	dictLock sync.Mutex
	//Cluster  *State
	addr net.Addr
	port uint16
}

func NewServer() *server {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil
	}
	ln, err := net.Listen("tcp", "127.0.0.1:18888")
	if err != nil {
		return nil
	}
	s := &server{
		logger:   logger,
		listener: ln,
		dict:     map[string]string{},
		dictLock: sync.Mutex{},
	}

	return s
}

func (s *server) ListenAndServe() {
	for {
		conn, e := s.listener.Accept()
		if e != nil {

		}

		go s.serveConn(conn)
	}
}

func (s *server) serveConn(conn net.Conn) {
	s.logger.Info(fmt.Sprintf("TCP: new client(%s)", conn.RemoteAddr()))
	s.processInlineBuffer(conn)
}

var separatorBytes = []byte(" ")

func (s *server) processInlineBuffer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {

		line, err := reader.ReadBytes('\n')
		if err != nil {
			return
		}
		length := len(line)
		if length <= 2 || line[length-2] != '\r' {
			// there are some empty lines within replication traffic, ignore this error
			//protocolError(ch, "empty line")
			continue
		}
		line = bytes.TrimSuffix(line, []byte{'\r', '\n'})
		fmt.Println(string(line))

		params := bytes.Split(line, separatorBytes)
		s.logger.Info(fmt.Sprintf("params:%s", params))
		switch {
		case bytes.Equal(params[0], []byte("SET")):
			key := string(params[1])
			value := string(params[2])
			s.SetKey(key, value)
			_, err := conn.Write([]byte("OK"))
			if err != nil {
				return
			}
		case bytes.Equal(params[0], []byte("GET")):
			key := string(params[1])
			v := s.GetValue(key)
			_, err := conn.Write([]byte(v))
			if err != nil {
				return
			}
		default:
			conn.Write([]byte("invalid command"))
		}
	}
}

func (s *server) SetKey(key, value string) {
	s.dictLock.Lock()
	defer s.dictLock.Unlock()
	s.dict[key] = value
}

func (s *server) GetValue(key string) string {
	s.dictLock.Lock()
	defer s.dictLock.Unlock()
	return s.dict[key]
}
