package factory

import (
	"chickChirick/internal/controller/abstraction"
	internalService "chickChirick/internal/controller/service/message"
	"net/http"
)

func InitMessageServer(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) {
	initDeletedService(mux, abstractDiContainer)
	initFileService(mux, abstractDiContainer)
	initGroupService(mux, abstractDiContainer)
	initMessageService(mux, abstractDiContainer)
	initMessageMetaService(mux, abstractDiContainer)
	initPersonalService(mux, abstractDiContainer)
	initSettingsService(mux, abstractDiContainer)
	initStatusService(mux, abstractDiContainer)
	initUserRelationService(mux, abstractDiContainer)
}

func initDeletedService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.DeletedController {
	deletedService := internalService.DeletedController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	deletedService.HandleRequest(mux)

	return deletedService
}

func initFileService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.FileController {
	fileService := internalService.FileController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	fileService.HandleRequest(mux)

	return fileService
}

func initGroupService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.GroupController {
	groupService := internalService.GroupController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	groupService.HandleRequest(mux)

	return groupService
}

func initMessageService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.MessagesController {
	messageService := internalService.MessagesController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	messageService.HandleRequest(mux)

	return messageService
}

func initMessageMetaService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.MetaController {
	metaService := internalService.MetaController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	metaService.HandleRequest(mux)

	return metaService
}

func initPersonalService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.PersonalController {
	personalService := internalService.PersonalController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	personalService.HandleRequest(mux)

	return personalService
}

func initSettingsService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.SettingsController {
	settingsService := internalService.SettingsController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	settingsService.HandleRequest(mux)

	return settingsService
}

func initStatusService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.StatusController {
	statusService := internalService.StatusController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	statusService.HandleRequest(mux)

	return statusService
}

func initUserRelationService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.UserRelationController {
	userRelationService := internalService.UserRelationController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	userRelationService.HandleRequest(mux)

	return userRelationService
}
