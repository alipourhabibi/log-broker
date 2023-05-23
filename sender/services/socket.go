package services

import (
	"fmt"
	"net"

	"github.com/alipourhabibi/log-broker/sender/configs"
)

type socket struct {
	conn net.Conn
}

func NewSocket() (*socket, error) {
	conn, err := net.Dial("tcp", configs.Confs.DestinationHost+":"+fmt.Sprintf("%d", configs.Confs.DestinationPort))
	if err != nil {
		return nil, err
	}
	return &socket{
		conn: conn,
	}, nil
}

func (s *socket) Send(message string) {
	s.conn.Write([]byte(message))
}
