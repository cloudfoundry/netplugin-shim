package service

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Service struct {
	done      chan bool
	waitGroup *sync.WaitGroup
	handle    handler
	logger    io.Writer
}

type handler func(net.Conn) error

func New(h handler) *Service {
	s := &Service{
		done:      make(chan bool),
		waitGroup: new(sync.WaitGroup),
		handle:    h,
	}
	s.waitGroup.Add(1)
	return s
}

func (s *Service) WithLogger(logger io.Writer) *Service {
	s.logger = logger
	return s
}

func (s *Service) Serve(listener *net.UnixListener) {
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.done:
			err := listener.Close()
			if err != nil {
				s.log(fmt.Sprintf("failed to close listener: %v", err))
			}
			return
		default:
		}

		err := listener.SetDeadline(time.Now().Add(time.Millisecond * 200))
		if err != nil {
			s.log(fmt.Sprintf("failed to set deadline on listener: %v", err))
			continue
		}

		conn, err := listener.AcceptUnix()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			s.log(fmt.Sprintf("failed to accept unix connection: %v", err))
		}
		s.waitGroup.Add(1)
		go s.serve(conn)
	}
}

func (s *Service) Stop() {
	close(s.done)
	s.waitGroup.Wait()
}

func (s *Service) serve(conn *net.UnixConn) {
	defer conn.Close()
	defer s.waitGroup.Done()

	//TODO: pick sensible value for this
	err := conn.SetDeadline(time.Now().Add(time.Millisecond * 800))
	if err != nil {
		s.log(fmt.Sprintf("failed to set deadline on connection: %v", err))
		return
	}

	if err := s.handle(conn); err != nil {
		s.log(fmt.Sprintf("failed to handle: %v", err))
	}
}

func (s *Service) log(message string) {
	if s.logger != nil {
		io.WriteString(s.logger, message)
	}
}
