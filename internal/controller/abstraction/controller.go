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

type RequestHandler interface {
	HandleRequest()
}

type ControllerInterface interface {
	RequestHandler
	Router()
}

//TODO: подумать как реализовать или удалить
//type RouteContainer struct {
//	DomainRoute string
//	Routes      []Route
//}
//
//type Route struct {
//	Name       string
//	UrlPattern *regexp.Regexp
//}
