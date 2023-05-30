package queue

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/alipourhabibi/log-broker/broker/configs"
	"github.com/alipourhabibi/log-broker/broker/internal/core/services/redis"
	"go.uber.org/zap"
)

type Server struct {
	RedisClient *redis.MessageQueue
	Client      net.Conn
	Server      net.Listener
	Logger      *zap.Logger
}
type ServerConfiguration func(*Server) error

func New(cfgs ...ServerConfiguration) (*Server, error) {
	s := &Server{}

	for _, cfg := range cfgs {
		err := cfg(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func WithZapper() ServerConfiguration {
	return func(s *Server) error {
		zap, err := zap.NewProduction()
		if err != nil {
			return err
		}
		s.Logger = zap
		return nil
	}
}

func WithRedisClient() ServerConfiguration {
	return func(s *Server) error {
		redisClient, err := redis.NewRedisClient()
		if err != nil {
			return err
		}
		s.RedisClient = redisClient
		return nil
	}
}

func WithClient() (net.Conn, error) {
	client, err := net.Dial("tcp", configs.Confs.DestinationHost+":"+fmt.Sprintf("%d", configs.Confs.DestinationPort))
	if err != nil {
		return nil, err
	}
	err = client.(*net.TCPConn).SetKeepAlive(true)
	if err != nil {
		return nil, err
	}
	err = client.(*net.TCPConn).SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func WithServer() ServerConfiguration {
	return func(s *Server) error {
		server, err := net.Listen("tcp", fmt.Sprintf(":%d", configs.Confs.Port))
		if err != nil {
			return err
		}
		s.Server = server
		return nil
	}
}

func (s *Server) Launch() {
	defer s.Logger.Sync()
	droped := make(chan bool, 1)
	connected := false
	s.Logger.Info("[INFO] Connecting to destination")
	client, err := WithClient()
	if err != nil {
		s.Logger.Error("[ERROR] Connection Failed to destination")
		droped <- true
	} else {
		s.Logger.Info("[INFO] Connection Succeed to destination")
		s.Client = client
		connected = true
	}
	go func() {
		for {
			select {
			case <-droped:
				s.Logger.Info("[INFO] Connecting to destination")
				client, err := WithClient()
				if err != nil {
					s.Logger.Error("[ERROR] Connection Failed to destination")
					droped <- true
					connected = false
					continue
				}
				s.Logger.Info("[INFO] Connection Succeed to destination")
				s.Client = client
				connected = true
			}
		}
	}()

	defer s.Server.Close()

	go func() {
		for {
			if connected {
				m, err := s.RedisClient.PopMessageFromQueue("logs")
				if err != nil {
					s.Logger.Error(fmt.Sprintf("[ERROR] %s", err.Error()))
					s.Logger.Info(fmt.Sprintf("[INFO] Pushing back %s", m))
					go s.RedisClient.PushMessageToQueue("logs", m)
					droped <- true
				}
				s.Logger.Info(fmt.Sprintf("[INFO] Sending %s", m))
				s.Client.Write([]byte(m + "\n"))
			} else {
				continue
			}
		}
	}()

	for {
		c, err := s.Server.Accept()
		if err != nil {
			s.Logger.Error("[ERROR] " + err.Error())
			continue
		}
		go s.handleConnection(c)
	}

}

func (s *Server) handleConnection(conn net.Conn) {
	for {
		buffer, _, err := bufio.NewReader(conn).ReadLine()
		if err != nil {
			conn.Close()
			return
		}
		err = s.RedisClient.PushMessageToQueue("logs", string(buffer))
		if err != nil {
			s.Logger.Error(fmt.Sprintf("[ERROR] %s", err.Error()))
		}
		s.Logger.Info(fmt.Sprintf("[INFO] Adding to the queue %s", string(buffer)))
	}
}
