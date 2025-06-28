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
	initController(container DIContainer)
	HandleRequest()
	GETId() int
}

func (c *Controller) initController(diContainer DIContainer) {
	c.Dependencies = diContainer
}

// HandleRequest TODO: удалить, если не будет использоваться
func (c *Controller) HandleRequest() {}

func (c *Controller) GETId(w http.ResponseWriter, r *http.Request) int {
	w.Header().Set("Content-Type", "application/json")
	idVal := r.URL.Query().Get("id")
	if idVal == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
	}
	id, _ := strconv.Atoi(idVal)

	return id
}
