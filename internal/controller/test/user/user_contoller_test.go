package u_test

import (
	"bytes"
	appCongig "chickChirick/cmd/config"
	"chickChirick/cmd/factory"
	"chickChirick/cmd/service"
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/service/user"
	"chickChirick/internal/middleware/config"
	userMiddleware "chickChirick/internal/middleware/validators/user"
	userModels "chickChirick/internal/model/user"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

//TODO: убедиться, что многократное открытие клиентов не привидет к ошибкам

// TODO: удалить хардкод
const ServerURL = "http://127.0.0.1:8081"

type ExternalServices struct {
	service.DBDecorator
	service.RedisDecorator
}

type UCContainer struct {
	ExternalServices
	HttpServer *httptest.Server
	HttpClient *http.Client
	*user.UserController
}

func initUCContainer(t *testing.T) UCContainer {
	factory.InitViper()

	appConfig := appCongig.AppConfiguration{}.NewAppConfiguration()

	dbDecorator := service.InitORM(&appConfig.DatabaseConfig)
	defer dbDecorator.CloseDB()

	redisDecorator := service.InitRedis(appConfig.RedisConfig)
	defer redisDecorator.RedisClose()

	ac := abstraction.Controller{
		Dependencies: abstraction.DIContainer{
			DBDecorator: dbDecorator,
		},
	}
	uv := userMiddleware.UserValidator{}

	//создание таблицы пользователя
	createUserTable(dbDecorator)

	//старт тестового сервера
	httpServer := startTestServer(t, dbDecorator, redisDecorator)

	return UCContainer{
		ExternalServices: ExternalServices{
			DBDecorator:    dbDecorator,
			RedisDecorator: redisDecorator,
		},
		HttpServer: httpServer,
		HttpClient: factory.InitHttpClient(),
		UserController: &user.UserController{
			Controller:    ac,
			UserValidator: uv,
		},
	}
}

func createUserTable(dbDecorator service.DBDecorator) {
	err := dbDecorator.GDB().AutoMigrate(userModels.User{})
	if err != nil {
		panic(err)
	}
}

func startTestServer(t *testing.T, db service.DBDecorator, redis service.RedisDecorator) *httptest.Server {
	t.Helper()

	mux := factory.BuildServer(db, redis)
	server := httptest.NewServer(mux)

	t.Cleanup(func() {
		server.Close()
	})

	return server
}

func TestCreateSuccess(t *testing.T) {
	ucc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	body, _ := json.Marshal(newUser)

	req, _ := http.NewRequest(http.MethodPost, ServerURL+"/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// middleware обходит контекст, эмуляция
	ctx := context.WithValue(req.Context(), config.UserUserKey, &newUser)
	req = req.WithContext(ctx)

	resp, err := ucc.HttpClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCreateAndGetSuccess(t *testing.T) {
	ucc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	body, _ := json.Marshal(newUser)

	req, _ := http.NewRequest(http.MethodPost, ServerURL+"/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// middleware обходит контекст, эмуляция
	ctx := context.WithValue(req.Context(), config.UserUserKey, &newUser)
	req = req.WithContext(ctx)

	resp, err := ucc.HttpClient.Do(req)

	//TODO: обработать ошибку
	defer resp.Body.Close()
	var j map[string]any
	err = json.NewDecoder(resp.Body).Decode(&j)

	//Получение списка всех пользователей
	//TODO: скорее всего не будет работать
	resp, err = ucc.HttpClient.Get(ServerURL + "/user?id=" + j["id"].(string))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var users []userModels.User
	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "Andrey", users[0].Name)
}
