package redis

import (
	"github.com/alipourhabibi/log-broker/broker/configs"
	"github.com/go-redis/redis"
)

type MessageQueue struct {
	Client *redis.Client
}

func NewRedisClient() (*MessageQueue, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     configs.Confs.RMQ.Host + ":" + configs.Confs.RMQ.Port, // Replace with your Redis server address
		Password: configs.Confs.RMQ.Password,                            // Set password if your Redis server requires authenticatio
		DB:       configs.Confs.RMQ.DBName,                              // Set the desired Redis database
	})

	// Ping the Redis server to check if the connection is successful
	_, err := rdb.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &MessageQueue{
		Client: rdb,
	}, nil
}

func (m *MessageQueue) PushMessageToQueue(queueName string, message string) error {
	err := m.Client.LPush(queueName, message).Err()
	if err != nil {
		return err
	}
	return nil
}

func (m *MessageQueue) PopMessageFromQueue(queueName string) (string, error) {
	result, err := m.Client.BRPop(0, queueName).Result()
	if err != nil {
		return "", err
	}
	message := result[1]
	return message, nil
}
