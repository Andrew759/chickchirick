package service

import (
	"chickChirick/cmd/configuration"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type RedisDecorator struct {
	Client *redis.Client
}

func InitAndPrepareRedis(config configuration.RedisConfig) RedisDecorator {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Username: config.User,
		Password: config.Password,
	})

	redisClient := RedisDecorator{
		Client: client,
	}

	redisClient.DeferRedisClose()

	return redisClient
}

func (rd RedisDecorator) DeferRedisClose() {
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			panic(fmt.Errorf("redis close error: %w", err))
		}
	}(rd.Client)
}
