package main

import (
	"chickChirick/cmd/configuration"
	"chickChirick/cmd/factory"
	"chickChirick/cmd/service"
	"fmt"
)

func main() {
	factory.InitViper()
	appConfig := configuration.NewConfiguration()

	dbDecorator := service.InitAndPrepareORM(appConfig.DatabaseConfig)

	redis := factory.InitRedis(appConfig.RedisConfig)
	httpProvider := factory.InitHttpClient()
	httpServer := factory.InitServer()

	//TODO: временная строка
	fmt.Println(dbDecorator, redis, httpProvider, httpServer)
}
