package u_test

import (
	"bytes"
	appConfig "chickChirick/cmd/config"
	"chickChirick/cmd/factory"
	"chickChirick/cmd/service"
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/http_transaction"
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
	"strconv"
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

	body, _ := json.Marshal(newUser)
	req, _ := http.NewRequest(http.MethodPost, ucc.ServerURL+"/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), config.UserUserKey, &newUser)
	req = req.WithContext(ctx)

	resp, _ := ucc.HttpClient.Do(req)
	defer resp.Body.Close()

	var result http_transaction.Response
	json.NewDecoder(resp.Body).Decode(&result)

	var createdUser userModels.User
	json.NewDecoder(result.ResultContainer).Decode(&createdUser)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, newUser.Name, createdUser.Name)
	assert.Equal(t, newUser.Surname, createdUser.Surname)
	assert.Equal(t, newUser.Phone, createdUser.Phone)
	assert.Equal(t, newUser.Login, createdUser.Login)
	assert.Equal(t, createdUser.CreatedAt, createdUser.UpdatedAt)
}

func TestCreateRepeatLoginFail(t *testing.T) {
	ucc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}

	body, _ := json.Marshal(newUser)
	req, _ := http.NewRequest(http.MethodPost, ucc.ServerURL+"/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), config.UserUserKey, &newUser)
	req = req.WithContext(ctx)

	resp, _ := ucc.HttpClient.Do(req)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	secondNewUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823343",
		Login:   "andrey_velkov",
	}

	body, _ = json.Marshal(secondNewUser)
	req, _ = http.NewRequest(http.MethodPost, ucc.ServerURL+"/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx = context.WithValue(req.Context(), config.UserUserKey, &newUser)
	req = req.WithContext(ctx)

	resp, _ = ucc.HttpClient.Do(req)
	defer resp.Body.Close()

	var result http_transaction.Response
	json.NewDecoder(resp.Body).Decode(&result)

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	assert.Equal(t, "user with login 'andrey_velkov' already exists", result.ErrorContainer.Message)
}

func TestCreateRepeatPhoneFail(t *testing.T) {
	ucc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}

	body, _ := json.Marshal(newUser)
	req, _ := http.NewRequest(http.MethodPost, ucc.ServerURL+"/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), config.UserUserKey, &newUser)
	req = req.WithContext(ctx)

	resp, _ := ucc.HttpClient.Do(req)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	secondNewUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov2",
	}

	body, _ = json.Marshal(secondNewUser)
	req, _ = http.NewRequest(http.MethodPost, ucc.ServerURL+"/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx = context.WithValue(req.Context(), config.UserUserKey, &newUser)
	req = req.WithContext(ctx)

	resp, _ = ucc.HttpClient.Do(req)
	defer resp.Body.Close()

	var result http_transaction.Response
	json.NewDecoder(resp.Body).Decode(&result)

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	assert.Equal(t, "user with phone '+79634823344' already exists", result.ErrorContainer.Message)
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
	req, _ := http.NewRequest(http.MethodPost, ucc.ServerURL+"/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), config.UserUserKey, &newUser)
	req = req.WithContext(ctx)

	resp, _ := ucc.HttpClient.Do(req)
	defer resp.Body.Close()

	var createResult http_transaction.Response
	json.NewDecoder(resp.Body).Decode(&createResult)
	var createdUser userModels.User
	json.NewDecoder(createResult.ResultContainer).Decode(&createdUser)

	resp, err := ucc.HttpClient.Get(ucc.ServerURL + "/user?id=" + strconv.Itoa(createdUser.Id))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var getResult http_transaction.Response
	json.NewDecoder(resp.Body).Decode(&getResult)
	var gotUser userModels.User
	json.NewDecoder(getResult.ResultContainer).Decode(&gotUser)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, newUser.Name, gotUser.Name)
	assert.Equal(t, newUser.Surname, gotUser.Surname)
	assert.Equal(t, newUser.Phone, gotUser.Phone)
	assert.Equal(t, newUser.Login, gotUser.Login)
	assert.Equal(t, createdUser.Id, gotUser.Id)
}

func TestCreateAndGetNotExistingUserFail(t *testing.T) {}

func TestCreateAndGetAll(t *testing.T) {}

func TestUpdateUserSuccess(t *testing.T) {}

func TestUpdateUserFail(t *testing.T) {}
