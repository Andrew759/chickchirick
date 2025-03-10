package factory

import (
	"fmt"
	"net/http"
)

// TODO: только для примера
type database map[string]int

func InitServer() *http.ServeMux {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	//две строчки выше можно упростить:
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	//log.Fatal(http.ListenAndServe("localhost:8000", mux))

	return mux
}

// TODO: только для примера
func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

// TODO: только для примера
func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}
