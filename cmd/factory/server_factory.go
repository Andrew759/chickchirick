package factory

import (
	"chickChirick/cmd/service"
	"chickChirick/internal/controller/c_controller"
	"log"
	"net/http"
	//TODO: подумать - оставить или удалить профилировщик
	_ "net/http/pprof"
)

func BuildAndServe(dbDecorator service.DBDecorator, redisDecorator service.RedisDecorator, httpClient *http.Client) {
	mux := BuildServer(dbDecorator, redisDecorator, httpClient)
	go func() {
		log.Println("Pprof server started on :6060")
		log.Fatal(http.ListenAndServe(":6060", nil))
	}()

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func BuildServer(dbDecorator service.DBDecorator, redisDecorator service.RedisDecorator, httpClient *http.Client) *http.ServeMux {
	mux := http.NewServeMux()

	container := c_controller.DIContainer{
		DBDecorator:    dbDecorator,
		RedisDecorator: redisDecorator,
	}

	InitUserServer(mux, container, httpClient)

	return mux
}
