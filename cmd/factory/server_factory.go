package factory

import (
	"net/http"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong"))
	if err != nil {
		panic(err)
	}
}

func InitServer() {
	//Регистрация обработчиков для конечных точек
	http.HandleFunc("/ping", pingHandler)

	// Запуск сервера
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
