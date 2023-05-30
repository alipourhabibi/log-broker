package services

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/alipourhabibi/log-broker/destination/configs"
)

type Server struct {
	Conn net.Listener
}

func NewServer() (*Server, error) {
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", configs.Confs.Port))
	if err != nil {
		return nil, err
	}
	return &Server{
		Conn: server,
	}, nil
}

func (s *Server) Launch() {
	defer s.Conn.Close()
	for {
		c, err := s.Conn.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go s.handleConnection(c)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	for {
		buffer, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			conn.Close()
			return
		}
		log.Printf("[INFO] Mssage Size is %d\n", len(buffer))
	}
}
