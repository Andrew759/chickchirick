package c_controller

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

type RequestHandler interface {
	HandleRequest()
}

type ControllerInterface interface {
	RequestHandler
	Router()
}
