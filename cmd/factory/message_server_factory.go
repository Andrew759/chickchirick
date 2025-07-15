package factory

import (
	"chickChirick/internal/controller/abstraction"
	internalService "chickChirick/internal/controller/service/message"
)

func InitMessageServer(abstractDiContainer abstraction.DIContainer) {
	_ = initDeletedService(abstractDiContainer)
	_ = initFileService(abstractDiContainer)
	_ = initGroupService(abstractDiContainer)
	_ = initMessageService(abstractDiContainer)
	_ = initMessageMetaService(abstractDiContainer)
	_ = initPersonalService(abstractDiContainer)
	_ = initSettingsService(abstractDiContainer)
	_ = initStatusService(abstractDiContainer)
	_ = initUserRelationService(abstractDiContainer)
}

func initDeletedService(abstractDiContainer abstraction.DIContainer) internalService.DeletedController {
	deletedService := internalService.DeletedController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	deletedService.HandleRequest()

	return deletedService
}

func initFileService(abstractDiContainer abstraction.DIContainer) internalService.FileController {
	fileService := internalService.FileController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	fileService.HandleRequest()

	return fileService
}

func initGroupService(abstractDiContainer abstraction.DIContainer) internalService.GroupController {
	groupService := internalService.GroupController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	groupService.HandleRequest()

	return groupService
}

func initMessageService(abstractDiContainer abstraction.DIContainer) internalService.MessagesController {
	messageService := internalService.MessagesController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}

	messageService.HandleRequest()

	return messageService
}

// InitMessageMetaService TODO: переименовать при разбиении на микросервисы
func initMessageMetaService(abstractDiContainer abstraction.DIContainer) internalService.MetaController {
	metaService := internalService.MetaController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	metaService.HandleRequest()

	return metaService
}

func initPersonalService(abstractDiContainer abstraction.DIContainer) internalService.PersonalController {
	personalService := internalService.PersonalController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	personalService.HandleRequest()

	return personalService
}

func initSettingsService(abstractDiContainer abstraction.DIContainer) internalService.SettingsController {
	settingsService := internalService.SettingsController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	settingsService.HandleRequest()

	return settingsService
}

func initStatusService(abstractDiContainer abstraction.DIContainer) internalService.StatusController {
	statusService := internalService.StatusController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	statusService.HandleRequest()

	return statusService
}

func initUserRelationService(abstractDiContainer abstraction.DIContainer) internalService.UserRelationController {
	userRelationService := internalService.UserRelationController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	userRelationService.HandleRequest()

	return userRelationService
}
