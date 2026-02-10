package factory

import (
	"chickChirick/internal/controller/c_controller"
	internalService "chickChirick/internal/controller/service/message"
	"net/http"
)

type MessageServer struct {
	*http.ServeMux
	c_controller.DIContainer
}

func InitMessageServer(mux *http.ServeMux, abstractDiContainer c_controller.DIContainer) MessageServer {
	messageServer := MessageServer{
		ServeMux:    mux,
		DIContainer: abstractDiContainer,
	}

	messageServer.initDeletedService()
	messageServer.initFileService()
	messageServer.initGroupService()
	messageServer.initMessageService()
	messageServer.initMessageMetaService()
	messageServer.initPersonalService()
	messageServer.initSettingsService()
	messageServer.initStatusService()
	messageServer.initUserRelationService()

	return messageServer
}

func (ms MessageServer) initDeletedService() internalService.DeletedController {
	deletedService := internalService.DeletedController{
		Controller: c_controller.Controller{
			ServeMux:     ms.ServeMux,
			Dependencies: ms.DIContainer,
		},
	}
	deletedService.HandleRequest()

	return deletedService
}

func (ms MessageServer) initFileService() internalService.FileController {
	fileService := internalService.FileController{
		Controller: c_controller.Controller{
			ServeMux:     ms.ServeMux,
			Dependencies: ms.DIContainer,
		},
	}
	fileService.HandleRequest()

	return fileService
}

func (ms MessageServer) initGroupService() internalService.GroupController {
	groupService := internalService.GroupController{
		Controller: c_controller.Controller{
			ServeMux:     ms.ServeMux,
			Dependencies: ms.DIContainer,
		},
	}
	groupService.HandleRequest()

	return groupService
}

func (ms MessageServer) initMessageService() internalService.MessagesController {
	messageService := internalService.MessagesController{
		Controller: c_controller.Controller{
			ServeMux:     ms.ServeMux,
			Dependencies: ms.DIContainer,
		},
	}
	messageService.HandleRequest()

	return messageService
}

func (ms MessageServer) initMessageMetaService() internalService.MetaController {
	metaService := internalService.MetaController{
		Controller: c_controller.Controller{
			ServeMux:     ms.ServeMux,
			Dependencies: ms.DIContainer,
		},
	}
	metaService.HandleRequest()

	return metaService
}

func (ms MessageServer) initPersonalService() internalService.PersonalController {
	personalService := internalService.PersonalController{
		Controller: c_controller.Controller{
			ServeMux:     ms.ServeMux,
			Dependencies: ms.DIContainer,
		},
	}
	personalService.HandleRequest()

	return personalService
}

func (ms MessageServer) initSettingsService() internalService.SettingsController {
	settingsService := internalService.SettingsController{
		Controller: c_controller.Controller{
			ServeMux:     ms.ServeMux,
			Dependencies: ms.DIContainer,
		},
	}
	settingsService.HandleRequest()

	return settingsService
}

func (ms MessageServer) initStatusService() internalService.StatusController {
	statusService := internalService.StatusController{
		Controller: c_controller.Controller{
			ServeMux:     ms.ServeMux,
			Dependencies: ms.DIContainer,
		},
	}
	statusService.HandleRequest()

	return statusService
}

func (ms MessageServer) initUserRelationService() internalService.UserRelationController {
	userRelationService := internalService.UserRelationController{
		Controller: c_controller.Controller{
			ServeMux:     ms.ServeMux,
			Dependencies: ms.DIContainer,
		},
	}
	userRelationService.HandleRequest()

	return userRelationService
}
