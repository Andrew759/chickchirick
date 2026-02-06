package factory

import (
	"chickChirick/internal/controller/abstraction"
	internalService "chickChirick/internal/controller/service/user"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/validators/user"
	"net/http"
)

type UserServer struct {
	*http.ServeMux
	abstraction.DIContainer
}

func InitUserServer(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) UserServer {
	userServer := UserServer{
		ServeMux:    mux,
		DIContainer: abstractDiContainer,
	}

	userServer.initUserService()
	userServer.initBanService()
	userServer.initUserMetaService()
	userServer.initPhotoService()
	userServer.initPropertyService()

	return userServer
}

func (us *UserServer) initUserService() internalService.UserController {
	userService := internalService.UserController{
		Controller: abstraction.Controller{
			ServeMux:     us.ServeMux,
			Dependencies: us.DIContainer,
		},
		Validator: user.UserValidatorFactory{}.NewValidator(us.DIContainer.DBDecorator),
	}
	userService.HandleRequest()

	return userService
}

func (us *UserServer) initBanService() internalService.BanController {
	banService := internalService.BanController{
		Controller: abstraction.Controller{
			ServeMux:     us.ServeMux,
			Dependencies: us.DIContainer,
		},
	}
	banService.HandleRequest()

	return banService
}

func (us *UserServer) initUserMetaService() internalService.MetaController {
	metaService := internalService.MetaController{
		Controller: abstraction.Controller{
			ServeMux:     us.ServeMux,
			Dependencies: us.DIContainer,
		},
	}
	metaService.HandleRequest()

	return metaService
}

func (us *UserServer) initPhotoService() internalService.PhotoController {
	photoService := internalService.PhotoController{
		Controller: abstraction.Controller{
			ServeMux:     us.ServeMux,
			Dependencies: us.DIContainer,
		},
	}
	photoService.HandleRequest()

	return photoService
}

func (us *UserServer) initPropertyService() internalService.PropertyController {
	propertyService := internalService.PropertyController{
		Controller: abstraction.Controller{
			ServeMux:     us.ServeMux,
			Dependencies: us.DIContainer,
		},
		Validator: user.PropertyValidatorFactory{}.NewValidator(us.DIContainer.DBDecorator, middleware.ActivateDBValidation()),
	}
	propertyService.HandleRequest()

	return propertyService
}
