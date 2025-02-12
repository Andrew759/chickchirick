package factory

import (
	"chickChirick/cmd/configuration"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func InitRedis(config configuration.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Username: config.User,
		Password: config.Password,
	})
}
