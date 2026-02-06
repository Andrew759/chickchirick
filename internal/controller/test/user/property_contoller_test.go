package u_test

import (
	"bytes"
	appConfig "chickChirick/cmd/config"
	"chickChirick/cmd/factory"
	"chickChirick/cmd/service"
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/controller/service/user"
	testAbstraction "chickChirick/internal/controller/test/abstraction"
	"chickChirick/internal/middleware/config"
	userMiddleware "chickChirick/internal/middleware/validators/user"
	userModels "chickChirick/internal/model/user"
	"chickChirick/pkg/chirik_faker"
	"chickChirick/pkg/chirik_gorm_tweaks/schema"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type PropertyControllerTestContainer struct {
	testAbstraction.ControllerTestContainer
	*user.PropertyController
}

func initPCContainer(t *testing.T) PropertyControllerTestContainer {
	t.Helper()
	factory.InitViper()

	appCfg := appConfig.AppConfiguration{}.NewAppConfiguration()
	db := service.InitORM(&appCfg.DatabaseConfig)
	redis := service.InitRedis(appCfg.RedisConfig)

	CreateUserTable(db)
	CreatePropertyTable(db)

	server := testAbstraction.StartTestServer(t, db, redis)

	ac := abstraction.Controller{
		Dependencies: abstraction.DIContainer{
			DBDecorator: db,
		},
	}
	pv := userMiddleware.PropertyValidator{}

	pCTC := PropertyControllerTestContainer{
		ControllerTestContainer: testAbstraction.ControllerTestContainer{
			ServerURL: appCfg.ServerURL,
			ExternalServices: testAbstraction.ExternalServices{
				DBDecorator:    db,
				RedisDecorator: redis,
			},
			HttpServer: server,
			HttpClient: factory.InitHttpClient(),
		},
		PropertyController: &user.PropertyController{
			Controller: ac,
			Validator:  pv,
		},
	}

	t.Cleanup(func() {
		DropUserTable(db)
		DropPropertyTable(db)
		pCTC.HttpServer.Close()
		pCTC.Controller.Dependencies.DBDecorator.CloseDB()
		pCTC.RedisDecorator.RedisClose()
	})

	return pCTC
}

func CreatePropertyTable(db service.DBDecorator) {
	if err := db.GDB().AutoMigrate(&userModels.Property{}); err != nil {
		panic("failed to migrate properties table: " + err.Error())
	}
}

func DropPropertyTable(db service.DBDecorator) {
	tableName, err := schema.GetTableName(db.GDB(), userModels.Property{})
	if err != nil {
		panic("failed to get table name: " + err.Error())
	}
	_, err = db.NativeDB().Exec("DROP TABLE IF EXISTS " + tableName + " CASCADE;")
	if err != nil {
		panic("failed to drop schema: " + err.Error())
	}
}

func doCreatePropertyRequest(t *testing.T, pctc PropertyControllerTestContainer, newProperty userModels.Property) (
	*http.Response,
	c_http.Response,
) {
	t.Helper()

	body, _ := json.Marshal(newProperty)
	req, _ := http.NewRequest(http.MethodPost, pctc.ServerURL+"/property", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), config.UserPropertyKey, &newProperty)
	req = req.WithContext(ctx)

	resp, _ := pctc.HttpClient.Do(req)
	defer resp.Body.Close()

	var decodedResponse c_http.Response
	json.NewDecoder(resp.Body).Decode(&decodedResponse)

	return resp, decodedResponse
}

func doUpdatePropertyRequest(t *testing.T, pctc UserControllerTestContainer, userId uint64, updatingProperty userModels.Property) (
	*http.Response,
	c_http.Response,
) {
	t.Helper()

	body, _ := json.Marshal(updatingProperty)
	req, _ := http.NewRequest(http.MethodPut, pctc.ServerURL+"/property/"+strconv.FormatUint(userId, 10), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := pctc.HttpClient.Do(req)
	defer resp.Body.Close()

	var decodedResponse c_http.Response
	json.NewDecoder(resp.Body).Decode(&decodedResponse)

	return resp, decodedResponse
}

func TestCreatePropertySuccess(t *testing.T) {
	pctc := initPCContainer(t)

	user := userModels.User{
		Id:      1,
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	userModels.CreateUser(pctc.GormInterface, &user)

	newProperty := userModels.Property{
		UserId:   user.Id,
		Timezone: 3,
		Email:    "Andreyvelkov@chirik.com",
	}
	newProperty.SetPassword("p@ssWoR_D1!")

	resp, decodedResponse := doCreatePropertyRequest(t, pctc, newProperty)

	var createdProperty userModels.Property
	json.NewDecoder(decodedResponse.PayloadContainer).Decode(&createdProperty)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, newProperty.Timezone, createdProperty.Timezone)
	assert.Equal(t, newProperty.Email, createdProperty.Email)
	assert.Equal(t, newProperty.Password, createdProperty.Password)
}

func TestCreatePropertyWithNotExistUserFail(t *testing.T) {
	pctc := initPCContainer(t)

	newProperty := userModels.Property{
		UserId:   9999,
		Timezone: 3,
		Email:    "Andreyvelkov@chirik.com",
	}
	newProperty.SetPassword("p@ssWoR_D1!")

	resp, decodedResp := doCreatePropertyRequest(t, pctc, newProperty)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, userModels.UserNotFoundErr, decodedResp.FirstError())
}

func TestCreatePropertyWithInvalidEmail(t *testing.T) {
	pctc := initPCContainer(t)

	user := userModels.User{
		Id:      1,
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	userModels.CreateUser(pctc.GormInterface, &user)

	newProperty := userModels.Property{
		UserId:   user.Id,
		Timezone: 3,
		Email:    "1_INVALID_EMAIL_@",
	}
	newProperty.SetPassword("p@ssWoR_D1!")

	resp, decodedResponse := doCreatePropertyRequest(t, pctc, newProperty)

	assert.Equal(t, "invalid email", decodedResponse.FirstError().Error())
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreatePropertyWithToLongPassword(t *testing.T) {
	pctc := initPCContainer(t)

	user := userModels.User{
		Id:      1,
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	userModels.CreateUser(pctc.GormInterface, &user)

	newProperty := userModels.Property{
		UserId:   user.Id,
		Timezone: 3,
		Email:    "Andreyvelkov@chirik.com",
	}
	newProperty.SetPassword(chirik_faker.FakeStringWithLength(1025))

	resp, decodedResponse := doCreatePropertyRequest(t, pctc, newProperty)

	assert.Equal(t, "invalid password", decodedResponse.FirstError().Error())
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreatePropertyIfItAlreadyExistFail(t *testing.T) {
	pctc := initPCContainer(t)

	user := userModels.User{
		Id:      1,
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	userModels.CreateUser(pctc.GormInterface, &user)

	firstProperty := userModels.Property{
		UserId:   user.Id,
		Timezone: 3,
		Email:    "Andreyvelkov@chirik.com",
	}
	firstProperty.SetPassword("p@ssWoR_D1!")
	resp, _ := doCreatePropertyRequest(t, pctc, firstProperty)

	secondProperty := userModels.Property{
		UserId:   user.Id,
		Timezone: 4,
		Email:    "Andreyvelkov2@chirik.com",
	}
	firstProperty.SetPassword("p@ssWoR_D2!")
	secondResp, secondDecodedResp := doCreatePropertyRequest(t, pctc, secondProperty)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, http.StatusConflict, secondResp.StatusCode)
	assert.Equal(t, userModels.PropertyForUserAlreadyExistsErr, secondDecodedResp.FirstError())
}

func TestCreateAndGetPropertySuccess(t *testing.T) {
	pctc := initPCContainer(t)

	user := userModels.User{
		Id:      1,
		Name:    "Andrey",
		Surname: "Velkov",
		Phone:   "+79634823344",
		Login:   "andrey_velkov",
	}
	userModels.CreateUser(pctc.GormInterface, &user)

	newProperty := userModels.Property{
		UserId:   1,
		Timezone: 3,
		Email:    "Andreyvelkov@chirik.com",
	}
	newProperty.SetPassword("p@ssWoR_D1!")

	resp, decodedResponse := doCreatePropertyRequest(t, pctc, newProperty)

	var createdProperty userModels.Property
	json.NewDecoder(decodedResponse.PayloadContainer).Decode(&createdProperty)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, newProperty.Timezone, createdProperty.Timezone)
	assert.Equal(t, newProperty.Email, createdProperty.Email)
	assert.Equal(t, newProperty.Password, createdProperty.Password)
}
