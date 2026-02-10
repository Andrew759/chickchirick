package factory

import (
	"chickChirick/cmd/service"
	"chickChirick/internal/controller/c_controller"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		panic(err)
	}
}

func BuildAndServe(dbDecorator service.DBDecorator, redisDecorator service.RedisDecorator) {
	mux := BuildServer(dbDecorator, redisDecorator)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func BuildServer(dbDecorator service.DBDecorator, redisDecorator service.RedisDecorator) *http.ServeMux {
	mux := http.NewServeMux()

	//Регистрация обработчиков для конечных точек
	mux.HandleFunc("/ping", pingHandler)

	container := c_controller.DIContainer{
		DBDecorator:    dbDecorator,
		RedisDecorator: redisDecorator,
	}

	InitUserServer(mux, container)
	InitMessageServer(mux, container)
	InitAuthServer(mux, container)

	return mux
}
