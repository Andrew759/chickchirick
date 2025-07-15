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

type ControllerInterface interface {
	HandleRequest()
}

// Init @deprecated
func (c *Controller) Init(container DIContainer) *Controller {
	c.Dependencies = container

	return c
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
