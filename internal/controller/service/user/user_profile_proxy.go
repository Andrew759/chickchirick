package user

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/dto"
	"chickChirick/internal/middleware/config"
	authValidator "chickChirick/internal/middleware/validator"
	"chickChirick/internal/model/user"
	"log/slog"
	"net/http"
)

type UserProfileProxy struct {
	Controller c_controller.Controller
	authValidator.AuthValidator
}

func (upp *UserProfileProxy) HandleRequest() {
	upp.Controller.ServeMux.HandleFunc("GET /frontend/profile", func(w http.ResponseWriter, r *http.Request) {
		upp.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			upp.GetProfile(w, c_http.NewRequest(r))
		})(w, r)
	})
}

func (upp *UserProfileProxy) GetProfile(w http.ResponseWriter, request *c_http.Request) {
	ctx := request.Context()

	uuidVal := ctx.Value(config.AuthKey)
	uuidStr, ok := uuidVal.(string)
	if !ok || uuidStr == "" {
		c_http.NewResponse().SendError(w, "Unauthorized: invalid user uuid in session", http.StatusUnauthorized)
		return
	}

	db := upp.Controller.Dependencies.DBDecorator.GDB().WithContext(ctx)

	u, err := user.GetUserByUuid(ctx, db, uuidStr)
	if err != nil {
		slog.Error("failed to get user by uuid: ", err.Error())
		c_http.NewResponse().SendError(w, "User profile not found", http.StatusNotFound)
		return
	}

	p, err := user.GetPropertyByUserId(ctx, db, u.Id)
	if err != nil {
		slog.Error("failed to get property by user id: ", err.Error())
		c_http.NewResponse().SendError(w, "Property not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w,
		dto.UserProfileResponse{
			Name:     u.Name,
			Surname:  u.Surname,
			Login:    u.Login,
			Phone:    u.Phone,
			Email:    p.Email,
			Password: p.Password,
		},
		http.StatusOK,
	)
}
