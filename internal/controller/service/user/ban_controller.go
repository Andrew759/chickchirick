package user

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	authValidator "chickChirick/internal/middleware/validator"
	ban "chickChirick/internal/model/user"
	"errors"
	"net/http"
)

type BanController struct {
	Controller c_controller.Controller
	middleware.Validator
	authValidator.AuthValidator
}

func (bc *BanController) HandleRequest() {
	bc.Controller.ServeMux.HandleFunc("GET /bans", func(w http.ResponseWriter, r *http.Request) {
		bc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			bc.GetBans(w, c_http.NewRequest(r))
		})
	})

	bc.Controller.ServeMux.HandleFunc("POST /ban", func(w http.ResponseWriter, r *http.Request) {
		bc.ValidateAuth(bc.Validate(func(w http.ResponseWriter, r *http.Request) {
			bc.CreateBan(w, c_http.NewRequest(r))
		}))
	})

	bc.Controller.ServeMux.HandleFunc("GET /ban/{id}", func(w http.ResponseWriter, r *http.Request) {
		bc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			bc.GetBan(w, c_http.NewRequest(r))
		})
	})

	bc.Controller.ServeMux.HandleFunc("GET /user/{id}/bans", func(w http.ResponseWriter, r *http.Request) {
		bc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			bc.GetBansByUserId(w, c_http.NewRequest(r))
		})
	})

	bc.Controller.ServeMux.HandleFunc("PUT /ban/{id}", func(w http.ResponseWriter, r *http.Request) {
		bc.ValidateAuth(bc.Validate(func(w http.ResponseWriter, r *http.Request) {
			bc.UpdateBan(w, c_http.NewRequest(r))
		}))
	})

	bc.Controller.ServeMux.HandleFunc("DELETE /ban/{id}", func(w http.ResponseWriter, r *http.Request) {
		bc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			bc.DeleteBan(w, c_http.NewRequest(r))
		})
	})
}

func (bc *BanController) GetBans(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()
	bans, err := ban.GetBans(ctx, bc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, bans, http.StatusOK)
}

func (bc *BanController) GetBan(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	b, err := ban.GetBanById(ctx, bc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Ban not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, b, http.StatusOK)
}

func (bc *BanController) GetBansByUserId(w http.ResponseWriter, r *c_http.Request) {
	userId, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	b, err := ban.GetBansByUserId(ctx, bc.Controller.Dependencies.DBDecorator.GDB(), userId)
	if err != nil {
		c_http.NewResponse().SendError(w, "Ban not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, b, http.StatusOK)
}

func (bc *BanController) CreateBan(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()
	b := ctx.Value(config.UserBanKey).(*ban.Ban)

	if err := ban.CreateBan(ctx, bc.Controller.Dependencies.DBDecorator.GDB(), b); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create ban: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, b, http.StatusCreated)
}

func (bc *BanController) UpdateBan(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	b := ctx.Value(config.UserBanKey).(*ban.Ban)

	err = ban.UpdateBanById(ctx, bc.Controller.Dependencies.DBDecorator.GDB(), b, id)
	if err != nil && errors.Is(err, ban.BanNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update ban: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, b, http.StatusOK)
}

func (bc *BanController) DeleteBan(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = ban.DeleteBanById(ctx, bc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, ban.BanNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete ban: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
