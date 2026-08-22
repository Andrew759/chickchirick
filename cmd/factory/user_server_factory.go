package factory

import (
	"chickChirick/internal/controller/c_controller"
	internalService "chickChirick/internal/controller/service/user"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/validator"
	"chickChirick/internal/middleware/validator/user"
	"chickChirick/pkg/chirik_config"
	"net/http"

	"github.com/spf13/viper"
)

type UserServer struct {
	*http.ServeMux
	c_controller.DIContainer

	//Кэш сервисов для проксирования
	UserService     internalService.UserController
	BanService      internalService.BanController
	UserMetaService internalService.MetaController
	PhotoService    internalService.PhotoController
	PropertyService internalService.PropertyController
}

func InitUserServer(mux *http.ServeMux, abstractDiContainer c_controller.DIContainer, httpClient *http.Client) UserServer {
	userServer := UserServer{
		ServeMux:    mux,
		DIContainer: abstractDiContainer,
	}

	userServer.UserService = userServer.initUserService()
	userServer.BanService = userServer.initBanService()
	userServer.UserMetaService = userServer.initUserMetaService()
	userServer.PhotoService = userServer.initPhotoService()
	userServer.PropertyService = userServer.initPropertyService()

	userServer.initUserProxyService()
	userServer.initUserProfileProxy()

	return userServer
}

func (us *UserServer) initUserService() internalService.UserController {
	userService := internalService.UserController{
		Controller: c_controller.Controller{
			ServeMux:     us.ServeMux,
			Dependencies: us.DIContainer,
		},
		Validator: user.UserValidatorFactory{}.NewValidator(us.DIContainer.DBDecorator),
		AuthValidator: validator.AuthValidator{
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
		AuthValidator: validator.AuthValidator{
			AuthServiceURL:   viper.GetString(chirik_config.AuthAppUrl),
			Client:           us.Client,
			ValidatorOptions: middleware.ValidatorOptions{},
		},
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
		AuthValidator: validator.AuthValidator{
			AuthServiceURL:   viper.GetString(chirik_config.AuthAppUrl),
			Client:           us.Client,
			ValidatorOptions: middleware.ValidatorOptions{},
		},
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
		AuthValidator: validator.AuthValidator{
			AuthServiceURL:   viper.GetString(chirik_config.AuthAppUrl),
			Client:           us.Client,
			ValidatorOptions: middleware.ValidatorOptions{},
		},
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
		AuthValidator: validator.AuthValidator{
			AuthServiceURL:   viper.GetString(chirik_config.AuthAppUrl),
			Client:           us.Client,
			ValidatorOptions: middleware.ValidatorOptions{},
		},
	}
	propertyService.HandleRequest()

	return propertyService
}

func (us *UserServer) initUserProxyService() internalService.CreateUserProxy {
	basicController := c_controller.Controller{
		ServeMux:     us.ServeMux,
		Dependencies: us.DIContainer,
	}

	userProxyService := internalService.CreateUserProxy{
		Controller: basicController,
		UPV:        user.UserProxyValidatorFactory{}.NewValidator(us.DIContainer.DBDecorator, middleware.ActivateDBValidation()),
		PV:         user.CreatePropertyValidatorFactory{}.NewValidator(us.DIContainer.DBDecorator),
	}
	userProxyService.HandleRequest()

	return userProxyService
}

func (us *UserServer) initUserProfileProxy() internalService.UserProfileProxy {
	profileService := internalService.UserProfileProxy{
		Controller: c_controller.Controller{
			ServeMux:     us.ServeMux,
			Dependencies: us.DIContainer,
		},
		AuthValidator: validator.AuthValidator{
			AuthServiceURL:   viper.GetString(chirik_config.AuthAppUrl),
			Client:           us.Client,
			ValidatorOptions: middleware.ValidatorOptions{},
		},
	}
	profileService.HandleRequest()

	return profileService
}
