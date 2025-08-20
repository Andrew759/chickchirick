package abstraction

import (
	"chickChirick/cmd/service"
	"net/http"
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
