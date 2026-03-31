package factory

import (
	"chickChirick/internal/controller/c_controller"
	internalService "chickChirick/internal/controller/service/user"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/validators"
	"chickChirick/internal/middleware/validators/user"
	"chickChirick/pkg/chirik_config"
	"net/http"

	"github.com/spf13/viper"
)

type UserServer struct {
	*http.ServeMux
	c_controller.DIContainer
	*http.Client
}

func InitUserServer(mux *http.ServeMux, abstractDiContainer c_controller.DIContainer, httpClient *http.Client) UserServer {
	userServer := UserServer{
		ServeMux:    mux,
		DIContainer: abstractDiContainer,
		Client:      httpClient,
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
		Controller: c_controller.Controller{
			ServeMux:     us.ServeMux,
			Dependencies: us.DIContainer,
		},
		Validator: user.UserValidatorFactory{}.NewValidator(us.DIContainer.DBDecorator),
		AuthValidator: validators.AuthValidator{
			AuthServiceURL:   viper.GetString(chirik_config.AuthAppUrl),
			Client:           us.Client,
			ValidatorOptions: middleware.ValidatorOptions{},
		},
	}
	userService.HandleRequest()

	return userService
}

func (us *UserServer) initBanService() internalService.BanController {
	banService := internalService.BanController{
		Controller: c_controller.Controller{
			ServeMux:     us.ServeMux,
			Dependencies: us.DIContainer,
		},
		Validator: user.BanValidatorFactory{}.NewValidator(us.DIContainer.DBDecorator, middleware.ActivateDBValidation()),
	}
	banService.HandleRequest()

	return banService
}

func (us *UserServer) initUserMetaService() internalService.MetaController {
	metaService := internalService.MetaController{
		Controller: c_controller.Controller{
			ServeMux:     us.ServeMux,
			Dependencies: us.DIContainer,
		},
		Validator: user.MetaValidatorFactory{}.NewValidator(us.DIContainer.DBDecorator, middleware.ActivateDBValidation()),
	}
	metaService.HandleRequest()

	return metaService
}

func (us *UserServer) initPhotoService() internalService.PhotoController {
	photoService := internalService.PhotoController{
		Controller: c_controller.Controller{
			ServeMux:     us.ServeMux,
			Dependencies: us.DIContainer,
		},
		Validator: user.PhotoValidatorFactory{}.NewValidator(us.DIContainer.DBDecorator, middleware.ActivateDBValidation()),
	}
	photoService.HandleRequest()

	return photoService
}

func (us *UserServer) initPropertyService() internalService.PropertyController {
	propertyService := internalService.PropertyController{
		Controller: c_controller.Controller{
			ServeMux:     us.ServeMux,
			Dependencies: us.DIContainer,
		},
		CPV: user.CreatePropertyValidatorFactory{}.NewValidator(us.DIContainer.DBDecorator, middleware.ActivateDBValidation()),
		UPV: user.UpdatePropertyValidatorFactory{}.NewValidator(us.DIContainer.DBDecorator, middleware.ActivateDBValidation()),
	}
	propertyService.HandleRequest()

	return propertyService
}
