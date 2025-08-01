package u_test

import (
	"bytes"
	appConfig "chickChirick/cmd/config"
	"chickChirick/cmd/factory"
	"chickChirick/cmd/service"
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/service/user"
	"chickChirick/internal/middleware/config"
	userMiddleware "chickChirick/internal/middleware/validators/user"
	userModels "chickChirick/internal/model/user"
	"chickChirick/pkg/chirik_gorm_tweaks"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ExternalServices struct {
	service.DBDecorator
	service.RedisDecorator
}

type UserControllerTestContainer struct {
	ServerURL string
	ExternalServices
	HttpServer *httptest.Server
	HttpClient *http.Client
	*user.UserController
}

func initUCContainer(t *testing.T) UserControllerTestContainer {
	t.Helper()
	factory.InitViper()

	appCfg := appConfig.AppConfiguration{}.NewAppConfiguration()
	db := service.InitORM(&appCfg.DatabaseConfig)
	redis := service.InitRedis(appCfg.RedisConfig)

	createUserTable(db)

	server := startTestServer(t, db, redis)

	ac := abstraction.Controller{
		Dependencies: abstraction.DIContainer{
			DBDecorator: db,
		},
	}
	uv := userMiddleware.UserValidator{}

	uCTC := UserControllerTestContainer{
		ServerURL: appCfg.ServerURL,
		ExternalServices: ExternalServices{
			DBDecorator:    db,
			RedisDecorator: redis,
		},
		HttpServer: server,
		HttpClient: factory.InitHttpClient(),
		UserController: &user.UserController{
			Controller:    ac,
			UserValidator: uv,
		},
	}

	t.Cleanup(func() {
		dropUserTable(uCTC)
		uCTC.HttpServer.Close()
		uCTC.DBDecorator.CloseDB()
		uCTC.RedisDecorator.RedisClose()
	})

	return uCTC
}

func createUserTable(db service.DBDecorator) {
	if err := db.GDB().AutoMigrate(&userModels.User{}); err != nil {
		panic("failed to migrate user table: " + err.Error())
	}
}

func dropUserTable(uCTC UserControllerTestContainer) {
	tableName, err := chirik_gorm_tweaks.GetTableName(uCTC.GDB(), userModels.User{})
	if err != nil {
		panic("failed to get table name: " + err.Error())
	}
	_, err = uCTC.NativeDB().Exec("DROP TABLE IF EXISTS " + tableName + " CASCADE;")
	if err != nil {
		panic("failed to drop schema: " + err.Error())
	}
}

func startTestServer(t *testing.T, db service.DBDecorator, redis service.RedisDecorator) *httptest.Server {
	t.Helper()
	mux := factory.BuildServer(db, redis)
	return httptest.NewServer(mux)
}

func TestCreateSuccess(t *testing.T) {
	ucc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	body, err := json.Marshal(newUser)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, ucc.ServerURL+"/user", bytes.NewBuffer(body))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), config.UserUserKey, &newUser)
	req = req.WithContext(ctx)

	resp, err := ucc.HttpClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestCreateAndGetSuccess(t *testing.T) {
	ucc := initUCContainer(t)

	user := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	body, err := json.Marshal(user)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, ucc.ServerURL+"/user", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), config.UserUserKey, &user)
	req = req.WithContext(ctx)

	resp, err := ucc.HttpClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	//resp, err = ucc.HttpClient.Get(ucc.ServerURL + "/user?id=" + id)
	//assert.NoError(t, err)
	//defer resp.Body.Close()
	//
	//assert.Equal(t, http.StatusOK, resp.StatusCode)
	//
	//var users []userModels.User
	//err = json.NewDecoder(resp.Body).Decode(&users)
	//assert.NoError(t, err)
	//assert.Len(t, users, 1)
	//assert.Equal(t, "Andrey", users[0].Name)
}
