package factory

import (
	"chickChirick/internal/controller/abstraction"
	internalService "chickChirick/internal/controller/service/auth"
	"net/http"
)

func InitAuthServer(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) {
	initCodeService(mux, abstractDiContainer)
	initSessionService(mux, abstractDiContainer)
	initTokenService(mux, abstractDiContainer)
}

func initCodeService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.CodeController {
	codeService := internalService.CodeController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	codeService.HandleRequest(mux)

	return codeService
}

func initSessionService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.SessionController {
	sessionService := internalService.SessionController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	sessionService.HandleRequest(mux)

	return sessionService
}

func initTokenService(mux *http.ServeMux, abstractDiContainer abstraction.DIContainer) internalService.TokenController {
	tokenService := internalService.TokenController{
		Controller: abstraction.Controller{
			Dependencies: abstractDiContainer,
		},
	}
	tokenService.HandleRequest(mux)

	return tokenService
}
