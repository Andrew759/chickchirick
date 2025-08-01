package factory

import (
	"chickChirick/internal/controller/abstraction"
	internalService "chickChirick/internal/controller/service/user"
	"net/http"
)

func InitUserServer(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) {
	initUserService(mux, abstractDiContainer)
	initBanService(mux, abstractDiContainer)
	initUserMetaService(mux, abstractDiContainer)
	initPhotoService(mux, abstractDiContainer)
	initPropertyService(mux, abstractDiContainer)
}

func initUserService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.UserController {
	userService := internalService.UserController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	userService.HandleRequest(mux)
	return userService
}

func initBanService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.BanController {
	banService := internalService.BanController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	banService.HandleRequest(mux)
	return banService
}

func initUserMetaService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.MetaController {
	metaService := internalService.MetaController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	metaService.HandleRequest(mux)
	return metaService
}

func initPhotoService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.PhotoController {
	photoService := internalService.PhotoController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	photoService.HandleRequest(mux)
	return photoService
}

func initPropertyService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.PropertyController {
	propertyService := internalService.PropertyController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	propertyService.HandleRequest(mux)
	return propertyService
}
