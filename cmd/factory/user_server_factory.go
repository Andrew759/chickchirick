package factory

import (
	"chickChirick/internal/controller/abstraction"
	internalService "chickChirick/internal/controller/service/user"
)

func InitUserServer(abstractDiContainer abstraction.DIContainer) {
	_ = initUserService(abstractDiContainer)
	_ = initBanService(abstractDiContainer)
	_ = initUserMetaService(abstractDiContainer)
	_ = initPhotoService(abstractDiContainer)
	_ = initPropertyService(abstractDiContainer)
}

func initUserService(abstractDiContainer abstraction.DIContainer) internalService.UserController {
	userService := internalService.UserController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	userService.HandleRequest()

	return userService
}

func initBanService(abstractDiContainer abstraction.DIContainer) internalService.BanController {
	banService := internalService.BanController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	banService.HandleRequest()

	return banService
}

// InitUserMetaService TODO: переименовать при разбиении на микросервисы
func initUserMetaService(abstractDiContainer abstraction.DIContainer) internalService.MetaController {
	metaService := internalService.MetaController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	metaService.HandleRequest()

	return metaService
}

func initPhotoService(abstractDiContainer abstraction.DIContainer) internalService.PhotoController {
	photoService := internalService.PhotoController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	photoService.HandleRequest()

	return photoService
}

func initPropertyService(abstractDiContainer abstraction.DIContainer) internalService.PropertyController {
	propertyService := internalService.PropertyController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	propertyService.HandleRequest()

	return propertyService
}
