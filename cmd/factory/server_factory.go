package factory

import (
	mainService "chickChirick/cmd/service"
	"chickChirick/internal/controller/abstraction"
	internalService "chickChirick/internal/controller/service/user"
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

	_ = initUserService(abstractDiContainer)
	_ = initBanService(abstractDiContainer)

	// Запуск сервера
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func initUserService(abstractDiContainer abstraction.DIContainer) internalService.UserController {
	userService := internalService.UserController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	userService.HandleRequest()

	return userService
}

func initBanService(abstractDiContainer abstraction.DIContainer) internalService.BanController {
	banService := internalService.BanController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	banService.HandleRequest()

	return banService
}
