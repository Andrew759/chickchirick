package factory

import (
	"chickChirick/cmd/service"
	"chickChirick/internal/controller/c_controller"
	middleware "chickChirick/internal/middleware/cors"
	"chickChirick/pkg/chirik_config"
	"log"
	"net/http"
	//TODO: подумать - оставить или удалить профилировщик
	_ "net/http/pprof"

	"github.com/spf13/viper"
)

func BuildAndServe(dbDecorator service.DBDecorator, redisDecorator service.RedisDecorator, httpClient *http.Client) {
	mux := BuildServer(dbDecorator, redisDecorator, httpClient)
	go func() {
		log.Println("Pprof server started on :6060")
		log.Fatal(http.ListenAndServe(":6060", nil))
	}()

	frontendURL := viper.GetString(chirik_config.FrontendUrl)
	handlerWithCORS := middleware.CORS(frontendURL, mux)

	err := http.ListenAndServe(":8080", handlerWithCORS)
	if err != nil {
		panic(err)
	}
}

func BuildServer(dbDecorator service.DBDecorator, redisDecorator service.RedisDecorator, httpClient *http.Client) *http.ServeMux {
	mux := http.NewServeMux()

	container := c_controller.DIContainer{
		DBDecorator:    dbDecorator,
		RedisDecorator: redisDecorator,
		Client:         httpClient,
	}

	InitUserServer(mux, container, httpClient)

	return mux
}
