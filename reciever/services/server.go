package services

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/alipourhabibi/log-broker/reciever/configs"
)

type serverSocket struct {
	l      net.Listener
	client *clientSocket
}

func NewServerSocket() (*serverSocket, error) {
	l, err := net.Listen("tcp", ":"+fmt.Sprintf("%d", configs.Confs.Port))
	if err != nil {
		return nil, err
	}
	client, err := NewClientSocket()
	if err != nil {
		return nil, err
	}
	return &serverSocket{
		l:      l,
		client: client,
	}, nil
}

func (s *serverSocket) Listen() {
	defer s.l.Close()
	for {
		c, err := s.l.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go s.handleConnection(c)
	}
}

func (s *serverSocket) handleConnection(conn net.Conn) {
	for {
		buffer, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			conn.Close()
			return
		}
		go s.client.Send(string(buffer))
	}
}
