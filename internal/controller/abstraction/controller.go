package abstraction

import (
	"chickChirick/cmd/service"
	"net/http"
	"strconv"
)

type DIContainer struct {
	DBDecorator    service.DBDecorator
	RedisDecorator service.RedisDecorator
}

type Controller struct {
	Dependencies DIContainer
}

type MainControllerInterface interface {
	InitController(container DIContainer)
	HandleRequest()
	GETId(w http.ResponseWriter, r *http.Request) int
}

func (c *Controller) InitController(container DIContainer) {
	c.Dependencies = container
}

func (c *Controller) HandleRequest() {}

func (c *Controller) GETId(w http.ResponseWriter, r *http.Request) int {
	w.Header().Set("Content-Type", "application/json")
	idVal := r.URL.Query().Get("id")
	if idVal == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}
	id, _ := strconv.Atoi(idVal)

	return id
}
