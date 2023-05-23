package services

import (
	"fmt"
	"net"

	"github.com/alipourhabibi/log-broker/reciever/configs"
)

type clientSocket struct {
	conn net.Conn
}

func NewClientSocket() (*clientSocket, error) {
	conn, err := net.Dial("tcp", configs.Confs.DestinationHost+":"+fmt.Sprintf("%d", configs.Confs.DestinationPort))
	if err != nil {
		return nil, err
	}
	return &clientSocket{
		conn: conn,
	}, nil
}

func (c *clientSocket) Send(message string) {
	c.conn.Write([]byte(message))
}
