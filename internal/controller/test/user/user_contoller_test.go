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
	"chickChirick/pkg/chirik_faker"
	"chickChirick/pkg/chirik_gorm_tweaks"
	"context"
	"encoding/json"
	"fmt"
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

// TODO doCreateUserRequest удалить t *testing.T, если не будет использоваться
func doCreateUserRequest(t *testing.T, uctc UserControllerTestContainer, newUser userModels.User) (
	*http.Response,
	http_transaction.Response,
) {
	body, _ := json.Marshal(newUser)
	req, _ := http.NewRequest(http.MethodPost, uctc.ServerURL+"/user", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), config.UserUserKey, &newUser)
	req = req.WithContext(ctx)

	resp, _ := uctc.HttpClient.Do(req)
	defer resp.Body.Close()

	var result http_transaction.Response
	json.NewDecoder(resp.Body).Decode(&result)

	return resp, result
}

func TestCreateSuccess(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

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
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	firstUserWithSameLoginResp, _ := doCreateUserRequest(t, uctc, newUser)

	secondNewUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823343",
		Login:   "andrey_velkov",
	}
	secondUserWithSameLoginResp, secondUserResult := doCreateUserRequest(t, uctc, secondNewUser)

	assert.Equal(t, http.StatusCreated, firstUserWithSameLoginResp.StatusCode)
	assert.Equal(t, http.StatusConflict, secondUserWithSameLoginResp.StatusCode)
	assert.Equal(t, "user with login 'andrey_velkov' already exists", secondUserResult.FirstError().Error())
}

func TestCreateRepeatPhoneFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	firstUserWithSamePhoneResp, _ := doCreateUserRequest(t, uctc, newUser)

	secondNewUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov2",
	}
	secondUserWithSamePhoneResp, secondUserResult := doCreateUserRequest(t, uctc, secondNewUser)

	assert.Equal(t, http.StatusCreated, firstUserWithSamePhoneResp.StatusCode)
	assert.Equal(t, http.StatusConflict, secondUserWithSamePhoneResp.StatusCode)
	assert.Equal(t, "user with phone '+79634823344' already exists", secondUserResult.FirstError().Error())
}

func TestCreateWithEmptyNameFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid name", result.FirstError().Error())
}

func TestCreateWithToLongNameFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    chirik_faker.FakeStringWithLength(257),
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid name", result.FirstError().Error())
}

func TestCreateWithEmptySurnameFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid surname", result.FirstError().Error())
}

func TestCreateWithToLongSurnameFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: chirik_faker.FakeStringWithLength(257),
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid surname", result.FirstError().Error())
}

func TestCreateWithEmptyPhoneFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "",
		Login:   "andrey_velkov",
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid phone", result.FirstError().Error())
}

func TestCreateWithToLongPhoneFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "99999999999999999",
		Login:   "andrey_velkov",
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid phone", result.FirstError().Error())
}

func TestCreateWithToShortPhoneFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "95144",
		Login:   "andrey_velkov",
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid phone", result.FirstError().Error())
}

func TestCreateWithInvalidPhoneFormatFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344_",
		Login:   "andrey_velkov",
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid phone", result.FirstError().Error())
}

func TestCreateWithEmptyLoginFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "",
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid login", result.FirstError().Error())
}

func TestCreateWithToShortLoginFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "s",
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid login", result.FirstError().Error())
}

func TestCreateWithToLongLoginFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   chirik_faker.FakeStringWithLength(257),
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid login", result.FirstError().Error())
}

func TestCreateWithInvalidLoginFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "@ll_names_are_coo!_",
	}
	resp, result := doCreateUserRequest(t, uctc, newUser)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "invalid login", result.FirstError().Error())
}

func TestCreateAndGetSuccess(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	_, createResult := doCreateUserRequest(t, uctc, newUser)

	var createdUser userModels.User
	json.NewDecoder(createResult.ResultContainer).Decode(&createdUser)

	getResp, err := uctc.HttpClient.Get(uctc.ServerURL + "/user?id=" + strconv.Itoa(createdUser.Id))
	defer getResp.Body.Close()

	var getResult http_transaction.Response
	json.NewDecoder(getResp.Body).Decode(&getResult)
	var gotUser userModels.User
	json.NewDecoder(getResult.ResultContainer).Decode(&gotUser)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)
	assert.Equal(t, newUser.Name, gotUser.Name)
	assert.Equal(t, newUser.Surname, gotUser.Surname)
	assert.Equal(t, newUser.Phone, gotUser.Phone)
	assert.Equal(t, newUser.Login, gotUser.Login)
	assert.Equal(t, createdUser.Id, gotUser.Id)
}

func TestCreateAndGetNotExistingUserFail(t *testing.T) {
	uctc := initUCContainer(t)

	newUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	_, createResult := doCreateUserRequest(t, uctc, newUser)

	var createdUser userModels.User
	json.NewDecoder(createResult.ResultContainer).Decode(&createdUser)

	notExistUserId := createdUser.Id + 1
	getResp, err := uctc.HttpClient.Get(uctc.ServerURL + "/user?id=" + strconv.Itoa(notExistUserId))
	defer getResp.Body.Close()

	var getResult http_transaction.Response
	json.NewDecoder(getResp.Body).Decode(&getResult)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
	assert.Equal(t, userModels.UserNotFoundErr, getResult.FirstError())
}

func TestCreateTwoUsersAndGetAll(t *testing.T) {
	uctc := initUCContainer(t)

	firstUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	_, firstUserResult := doCreateUserRequest(t, uctc, firstUser)

	secondUser := userModels.User{
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823343",
		Login:   "andrey_velkov2",
	}
	_, secondUserResult := doCreateUserRequest(t, uctc, secondUser)

	getAllResp, err := uctc.HttpClient.Get(uctc.ServerURL + "/users")
	defer getAllResp.Body.Close()

	var getAllResult http_transaction.Response
	json.NewDecoder(getAllResp.Body).Decode(&getAllResult)

	assert.NoError(t, err)

	//TODO: удалить. Также доработать автотест. Сейчас он ничего не проверяет!
	fmt.Println(firstUserResult, secondUserResult)
}

func TestUpdateUserSuccess(t *testing.T) {}

func TestUpdateUserFail(t *testing.T) {}
