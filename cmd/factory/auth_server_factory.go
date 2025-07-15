package factory

import (
	"chickChirick/internal/controller/abstraction"
	internalService "chickChirick/internal/controller/service/auth"
)

func InitAuthServer(abstractDiContainer abstraction.DIContainer) {
	_ = initCodeService(abstractDiContainer)
	_ = initSessionService(abstractDiContainer)
	_ = initTokenService(abstractDiContainer)
}

func initCodeService(abstractDiContainer abstraction.DIContainer) internalService.CodeController {
	codeService := internalService.CodeController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	codeService.HandleRequest()

	return codeService
}

func initSessionService(abstractDiContainer abstraction.DIContainer) internalService.SessionController {
	sessionService := internalService.SessionController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	sessionService.HandleRequest()

	return sessionService
}

func initTokenService(abstractDiContainer abstraction.DIContainer) internalService.TokenController {
	tokenService := internalService.TokenController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	tokenService.HandleRequest()

	return tokenService
}
