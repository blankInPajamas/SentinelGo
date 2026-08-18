package syslog

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type SysLogCollector struct {
	mu sync.Mutex
	conn *net.UDPConn
	addr string
}

func New(addr string) *SysLogCollector {
	return &SysLogCollector{
		addr: addr,
	}
}

func (s *SysLogCollector) Start(ctx context.Context, handler func(line string)) error {
		addr, err := net.ResolveUDPAddr("udp", s.addr)

		if err != nil {
			return fmt.Errorf("resolving UDP addr: %w", err)
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			return fmt.Errorf("listening on UDP: %w", err)
		}

		s.mu.Lock()
		s.conn = conn
		s.mu.Unlock()

		defer func() {
			s.mu.Lock()
			if s.conn != nil {
				s.conn.Close()
			}
			s.mu.Unlock()
		}()
	

	buffer := make([]byte, 65535)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err() 
		default:

			s.conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				return fmt.Errorf("reading UDP: %w", err)
			}

			line := string(buffer[:n])
			handler(line)
		}
	}
}

func (s *SysLogCollector) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}