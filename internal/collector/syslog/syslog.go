package syslog

import (
	"net"
	"sync"
)

type SysLogCollector struct {
	mu sync.Mutex
	connec *net.UDPConn
	addr string
}