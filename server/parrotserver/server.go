package parrotserver

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"

	"go.uber.org/zap"
)

type server struct {
	lg       *zap.Logger
	errChan  chan error
	listener net.Listener
	dict     map[string]string
	dictLock sync.Mutex
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
		lg:       logger,
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

func (s server) serveConn(conn net.Conn) {
	s.lg.Info(fmt.Sprintf("TCP: new client(%s)", conn.RemoteAddr()))
	s.processInlineBuffer(conn)
}

var separatorBytes = []byte(" ")

func (s server) processInlineBuffer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = fmt.Errorf("failed to read command - %s", err)
			}
			return
		}

		// trim the '\n'
		line = line[:len(line)-1]
		// optionally trim the '\r'
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		params := bytes.Split(line, separatorBytes)
		s.lg.Info(fmt.Sprintf("params:%s", params))
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
