package services

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var logLevels = []string{"INFO", "DEBUG", "WARNING", "ERROR"}
var logSources = []string{"App", "System", "Network", "Security"}

type logService struct {
	socket *socket
}

func NewLogService() (*logService, error) {
	socket, err := NewSocket()
	if err != nil {
		return nil, err
	}
	return &logService{
		socket: socket,
	}, nil
}

func (l *logService) Run() error {
	wg := sync.WaitGroup{}
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	messages := make(chan string, 15000)

	wg.Add(1)

	// Genearte Messages
	go func() {
		for {
			select {
			case <-done:
				wg.Done()
				return
			case <-ticker.C:
				for i := 1; i <= 10000+rand.Intn(5000); i++ {
					messages <- l.generateRandomLogMessage()
				}
			}
		}
	}()

	// Send Messages
	go func() {
		for m := range messages {
			go l.socket.Send(m)
		}
	}()

	wg.Wait()
	return nil
}

func (l *logService) generateRandomLogMessage() string {
	rand.Seed(time.Now().Unix())
	curTime := time.Now().Format("2006-01-02 15:04:05")
	logLevel := logLevels[rand.Intn(len(logLevels))]
	logSource := logSources[rand.Intn(len(logSources))]
	logMessage := l.generateRandomString(rand.Intn(8*1000-50) + 50)
	logMessage = logMessage[len(curTime)+len(logLevel)+len(logSource):]

	logEntry := fmt.Sprintf("%s [%s] %s %s", curTime, logLevel, logSource, logMessage)
	return logEntry
}

func (l *logService) generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
