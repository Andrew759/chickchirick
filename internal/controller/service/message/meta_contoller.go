package message

import (
	"chickChirick/internal/controller/abstraction"
	meta "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

type MetaController struct {
	Controller abstraction.Controller
}

func (mc *MetaController) HandleRequest() {
	http.HandleFunc("/message/metas", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			mc.GetMetas(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/message/meta", func(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (mc *MetaController) GetMetas(w http.ResponseWriter) {
	metas, err := meta.GetMetas(mc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(metas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (mc *MetaController) GetMeta(w http.ResponseWriter, r *http.Request) {
	id := mc.Controller.GETId(w, r)
	m, err := meta.GetMetaById(mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Meta not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (mc *MetaController) CreateMeta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m meta.Meta
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := meta.CreateMeta(mc.Controller.Dependencies.DBDecorator.GDB(), &m); err != nil {
		http.Error(w, "Failed to create meta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (mc *MetaController) UpdateMeta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m meta.Meta
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := meta.UpdateMeta(mc.Controller.Dependencies.DBDecorator.GDB(), &m)
	if err != nil {
		http.Error(w, "Failed to update meta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (mc *MetaController) DeleteMeta(w http.ResponseWriter, r *http.Request) {
	id := mc.Controller.GETId(w, r)
	err := meta.DeleteMetaById(mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete meta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
