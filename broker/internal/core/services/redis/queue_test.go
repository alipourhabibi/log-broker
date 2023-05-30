package redis

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	m, err := NewRedisClient()
	if err != nil {
		t.Errorf(err.Error())
	}
	go func() {
		m.PushMessageToQueue("ali", "mamad")
		m.PushMessageToQueue("ali", "mamad")
		m.PushMessageToQueue("ali", "mamad")
		m.PopMessageFromQueue("ali")
	}()
	for {
		fmt.Println(m.PopMessageFromQueue("ali"))
	}
}
