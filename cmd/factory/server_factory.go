package factory

import (
	"log"
	"net/http"
)

func InitServer() *http.ServeMux {
	mux := http.NewServeMux()
	log.Fatal(http.ListenAndServe("localhost:8000", mux))

	return mux
}
