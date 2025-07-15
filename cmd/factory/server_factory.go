package factory

import (
	mainService "chickChirick/cmd/service"
	"chickChirick/internal/controller/abstraction"
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		panic(err)
	}
}

func InitServer(dbDecorator mainService.DBDecorator, redisDecorator mainService.RedisDecorator) {
	//Регистрация обработчиков для конечных точек
	http.HandleFunc("/ping", pingHandler)

	abstractDiContainer := abstraction.DIContainer{
		DBDecorator:    dbDecorator,
		RedisDecorator: redisDecorator,
	}

	InitUserServer(abstractDiContainer)
	InitMessageServer(abstractDiContainer)
	InitAuthServer(abstractDiContainer)

	// Запуск сервера
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
