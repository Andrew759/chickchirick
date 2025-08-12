package abstraction

import (
	"chickChirick/cmd/service"
	"fmt"
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
// @deprecated
// TODO: метод должен быть заменён на аналогичный из c_http.Request
func (c *Controller) HttpId(w http.ResponseWriter, r *http.Request, urlPrefix string) (int, error) {
	path := strings.TrimPrefix(r.URL.Path, urlPrefix)
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		return 0, fmt.Errorf("missing id in URL")
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		return 0, fmt.Errorf("invalid URL id")
	}

	return id, nil
}
