package user

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/middleware"
	"chickChirick/internal/middleware/config"
	authValidator "chickChirick/internal/middleware/validators"
	meta "chickChirick/internal/model/user"
	"errors"
	"net/http"
)

type MetaController struct {
	Controller c_controller.Controller
	middleware.Validator
	authValidator.AuthValidator
}

// HandleRequest TODO: здесь временно определяется путь user для будущего микросервиса, т.к существует пересечение с
// meta_controller в message
func (mc *MetaController) HandleRequest() {
	mc.Controller.ServeMux.HandleFunc("GET /metas", func(w http.ResponseWriter, r *http.Request) {
		mc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			mc.GetMetas(w, c_http.NewRequest(r))
		})
	})

	mc.Controller.ServeMux.HandleFunc("POST /meta",
		mc.Validate(func(w http.ResponseWriter, r *http.Request) {
			mc.CreateMeta(w, c_http.NewRequest(r))
		}))

	mc.Controller.ServeMux.HandleFunc("GET /user/{id}/meta", func(w http.ResponseWriter, r *http.Request) {
		mc.ValidateAuth(func(w http.ResponseWriter, r *http.Request) {
			mc.GetMeta(w, c_http.NewRequest(r))
		})
	})

	mc.Controller.ServeMux.HandleFunc("PUT /user/{id}/meta",
		mc.ValidateAuth(mc.Validate(func(w http.ResponseWriter, r *http.Request) {
			mc.UpdateMeta(w, c_http.NewRequest(r))
		})))

	mc.Controller.ServeMux.HandleFunc("DELETE /user/{id}/meta", func(w http.ResponseWriter, r *http.Request) {
		mc.ValidateAuth(mc.Validate(func(w http.ResponseWriter, r *http.Request) {
			mc.DeleteMeta(w, c_http.NewRequest(r))
		}))
	})
}

func (mc *MetaController) GetMetas(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()
	metas, err := meta.GetMetas(ctx, mc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, metas, http.StatusOK)
}

func (mc *MetaController) GetMeta(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	m, err := meta.GetMetaByUserId(ctx, mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Meta not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, m, http.StatusOK)
}

func (mc *MetaController) CreateMeta(w http.ResponseWriter, r *c_http.Request) {
	ctx := r.Context()
	m := ctx.Value(config.UserMetaKey).(*meta.Meta)

	if err := meta.CreateMeta(ctx, mc.Controller.Dependencies.DBDecorator.GDB(), m); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create meta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, m, http.StatusOK)
}

func (mc *MetaController) UpdateMeta(w http.ResponseWriter, r *c_http.Request) {
	userId, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	m := ctx.Value(config.UserMetaKey).(*meta.Meta)

	err = meta.UpdateMetaByUserId(ctx, mc.Controller.Dependencies.DBDecorator.GDB(), m, userId)
	if err != nil && errors.Is(err, meta.MetaNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update meta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, m, http.StatusOK)
}

func (mc *MetaController) DeleteMeta(w http.ResponseWriter, r *c_http.Request) {
	userId, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = meta.DeleteMetaByUserId(ctx, mc.Controller.Dependencies.DBDecorator.GDB(), userId)
	if err != nil && errors.Is(err, meta.MetaNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete meta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
