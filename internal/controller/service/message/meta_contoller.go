package message

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	meta "chickChirick/internal/model/message"
	"encoding/json"
	"errors"
	"net/http"
)

type MetaController struct {
	Controller c_controller.Controller
}

func (mc *MetaController) HandleRequest() {
	mc.Controller.ServeMux.HandleFunc("GET /metas", func(w http.ResponseWriter, r *http.Request) {
		mc.GetMetas(w)
	})

	mc.Controller.ServeMux.HandleFunc("POST /meta", func(w http.ResponseWriter, r *http.Request) {
		mc.CreateMeta(w, c_http.NewRequest(r))
	})

	mc.Controller.ServeMux.HandleFunc("GET /meta/{id}", func(w http.ResponseWriter, r *http.Request) {
		mc.GetMeta(w, c_http.NewRequest(r))
	})

	mc.Controller.ServeMux.HandleFunc("PUT /meta/{id}", func(w http.ResponseWriter, r *http.Request) {
		mc.UpdateMeta(w, c_http.NewRequest(r))
	})

	mc.Controller.ServeMux.HandleFunc("DELETE /meta/{id}", func(w http.ResponseWriter, r *http.Request) {
		mc.DeleteMeta(w, c_http.NewRequest(r))
	})
}

func (mc *MetaController) GetMetas(w http.ResponseWriter) {
	metas, err := meta.GetMetas(mc.Controller.Dependencies.DBDecorator.GDB())
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

	m, err := meta.GetMetaById(mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Meta not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, m, http.StatusOK)
}

func (mc *MetaController) CreateMeta(w http.ResponseWriter, r *c_http.Request) {
	var m meta.Meta
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := meta.CreateMeta(mc.Controller.Dependencies.DBDecorator.GDB(), &m); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create meta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, m, http.StatusCreated)
}

func (mc *MetaController) UpdateMeta(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var m meta.Meta
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = meta.UpdateMetaById(mc.Controller.Dependencies.DBDecorator.GDB(), &m, id)
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
	id, err := r.HTTPId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = meta.DeleteMetaById(mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, meta.MetaNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete meta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
