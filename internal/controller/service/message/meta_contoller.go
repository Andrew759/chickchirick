package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	meta "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

type MetaController struct {
	Controller abstraction.Controller
}

func (mc *MetaController) HandleRequest() {
	mc.Controller.ServeMux.HandleFunc("/message/metas", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			mc.GetMetas(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mc.Controller.ServeMux.HandleFunc(MetaResource, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			mc.GetMeta(w, r)
		case http.MethodPost:
			mc.CreateMeta(w, r)
		case http.MethodPut:
			mc.UpdateMeta(w, r)
		case http.MethodDelete:
			mc.DeleteMeta(w, r)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
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

func (mc *MetaController) GetMeta(w http.ResponseWriter, r *http.Request) {
	id, err := mc.Controller.HttpId(w, r, MetaResource)
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
	}

	m, err := meta.GetMetaById(mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Meta not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, m, http.StatusOK)
}

func (mc *MetaController) CreateMeta(w http.ResponseWriter, r *http.Request) {
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

func (mc *MetaController) UpdateMeta(w http.ResponseWriter, r *http.Request) {
	var m meta.Meta
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := meta.UpdateMeta(mc.Controller.Dependencies.DBDecorator.GDB(), &m)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update meta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, m, http.StatusOK)
}

func (mc *MetaController) DeleteMeta(w http.ResponseWriter, r *http.Request) {
	id, err := mc.Controller.HttpId(w, r, MetaResource)
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
	}

	err = meta.DeleteMetaById(mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete meta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
