package abstraction

import (
	"chickChirick/cmd/factory"
	"chickChirick/cmd/service"
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ControllerTestContainerInterface interface {
	abstraction.ControllerInterface
	middleware.Validator
}

type ExternalServices struct {
	service.DBDecorator
	service.RedisDecorator
}

type ControllerTestContainer struct {
	ServerURL string
	ExternalServices
	HttpServer *httptest.Server
	HttpClient *http.Client
}

func StartTestServer(t *testing.T, db service.DBDecorator, redis service.RedisDecorator) *httptest.Server {
	t.Helper()
	mux := factory.BuildServer(db, redis)
	return httptest.NewServer(mux)
}
