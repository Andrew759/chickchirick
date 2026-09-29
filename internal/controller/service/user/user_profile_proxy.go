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
	"strings"
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

	upp.Controller.ServeMux.HandleFunc("GET /frontend/users/by-uuids", func(w http.ResponseWriter, r *http.Request) {
		upp.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			upp.GetUsersByUuids(w, c_http.NewRequest(r))
		})(w, r)
	})

	//TODO: вынести в отдельный сервис
	upp.Controller.ServeMux.HandleFunc("GET /frontend/users/search", func(w http.ResponseWriter, r *http.Request) {
		upp.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			upp.SearchUsers(w, c_http.NewRequest(r))
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

func (upp *UserProfileProxy) GetUsersByUuids(w http.ResponseWriter, request *c_http.Request) {
	ctx := request.Context()

	raw := request.URL.Query().Get("uuids")
	if strings.TrimSpace(raw) == "" {
		c_http.NewResponse().SendSuccess(w, []dto.UserBrief{}, http.StatusOK)
		return
	}

	parts := strings.Split(raw, ",")
	uuids := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, p := range parts {
		u := strings.TrimSpace(p)
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		uuids = append(uuids, u)
	}

	if len(uuids) == 0 {
		c_http.NewResponse().SendSuccess(w, []dto.UserBrief{}, http.StatusOK)
		return
	}

	db := upp.Controller.Dependencies.DBDecorator.GDB()
	rows, err := user.GetUsersByUuids(ctx, db, uuids)
	if err != nil {
		slog.Error("GetUsersByUuids failed: ", err.Error())
		c_http.NewResponse().SendError(w, "failed to load users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	out := make([]dto.UserBrief, 0, len(rows))
	for _, row := range rows {
		out = append(out, dto.UserBrief{
			ID:       row.Id,
			UserUuid: row.UserUuid,
			Name:     row.Name,
			Surname:  row.Surname,
			Login:    row.Login,
		})
	}

	c_http.NewResponse().SendSuccess(w, out, http.StatusOK)
}

// SearchUsers ищет пользователей по логину или телефону (query param q).
// В ответе — id user-сервиса и userUuid (нужен для связи с messages).
func (upp *UserProfileProxy) SearchUsers(w http.ResponseWriter, request *c_http.Request) {
	ctx := request.Context()
	q := strings.TrimSpace(request.URL.Query().Get("q"))
	if q == "" {
		c_http.NewResponse().SendSuccess(w, []dto.UserBrief{}, http.StatusOK)
		return
	}
	if len(q) < 2 {
		c_http.NewResponse().SendError(w, "query must be at least 2 characters", http.StatusBadRequest)
		return
	}

	db := upp.Controller.Dependencies.DBDecorator.GDB()

	excludeUuid := ""
	if uuidVal := ctx.Value(config.AuthKey); uuidVal != nil {
		if uuidStr, ok := uuidVal.(string); ok {
			excludeUuid = uuidStr
		}
	}

	rows, err := user.SearchUsersWithUuid(ctx, db, q, excludeUuid, 20)
	if err != nil {
		slog.Error("SearchUsers failed: ", err.Error())
		c_http.NewResponse().SendError(w, "search failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	out := make([]dto.UserBrief, 0, len(rows))
	for _, row := range rows {
		out = append(out, dto.UserBrief{
			ID:       row.Id,
			UserUuid: row.UserUuid,
			Name:     row.Name,
			Surname:  row.Surname,
			Login:    row.Login,
		})
	}

	c_http.NewResponse().SendSuccess(w, out, http.StatusOK)
}
