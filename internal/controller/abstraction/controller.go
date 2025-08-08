package abstraction

import (
	"chickChirick/cmd/service"
	"chickChirick/internal/controller/http_transaction"
	"net/http"
	"strconv"
	"strings"
)

type DIContainer struct {
	DBDecorator    service.DBDecorator
	RedisDecorator service.RedisDecorator
}

type Controller struct {
	ServeMux     *http.ServeMux
	Dependencies DIContainer
}

type ControllerInterface interface {
	HandleRequest()
}

func (c *Controller) HandleRequest() {}

// HttpId - берет id из URL. Метод предполагает, что ID передается в согласовании с правилами REST API
func (c *Controller) HttpId(w http.ResponseWriter, r *http.Request, urlPrefix string) int {
	idStr := strings.TrimPrefix(r.URL.Path, urlPrefix)
	if idStr == "" || idStr == r.URL.Path {
		http_transaction.NewResponse().SendError(w, "Invalid URL", http.StatusBadRequest)
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http_transaction.NewResponse().SendError(w, "Invalid URL id", http.StatusBadRequest)
	}

	return id
}
