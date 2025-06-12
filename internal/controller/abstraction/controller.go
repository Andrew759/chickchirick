package abstraction

import "chickChirick/cmd/service"

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
}

func (c *Controller) initController(diContainer DIContainer) {
	c.Dependencies = diContainer
}
func (c *Controller) HandleRequest() {}
