package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	settings "chickChirick/internal/model/message"
	"encoding/json"
	"errors"
	"net/http"
)

type SettingsController struct {
	Controller abstraction.Controller
}

func (sc *SettingsController) HandleRequest() {
	sc.Controller.ServeMux.HandleFunc("/message/settings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sc.GetSettings(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	sc.Controller.ServeMux.HandleFunc("/message/setting", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			sc.CreateSetting(w, c_http.NewRequest(r))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	sc.Controller.ServeMux.HandleFunc("/message/setting/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sc.GetSetting(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/message/setting/")))
		case http.MethodPut:
			sc.UpdateSetting(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/message/setting/")))
		case http.MethodDelete:
			sc.DeleteSetting(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/message/setting/")))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (sc *SettingsController) GetSettings(w http.ResponseWriter) {
	settingList, err := settings.GetSettings(sc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, settingList, http.StatusOK)
}

func (sc *SettingsController) GetSetting(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	s, err := settings.GetSettingsById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendSuccess(w, "Setting not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, s, http.StatusOK)
}

func (sc *SettingsController) CreateSetting(w http.ResponseWriter, r *c_http.Request) {
	var s settings.Settings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		c_http.NewResponse().SendSuccess(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := settings.CreateSettings(sc.Controller.Dependencies.DBDecorator.GDB(), &s); err != nil {
		c_http.NewResponse().SendSuccess(w, "Failed to create setting: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, "Success", http.StatusCreated)
}

func (sc *SettingsController) UpdateSetting(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var s settings.Settings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = settings.UpdateSettingsById(sc.Controller.Dependencies.DBDecorator.GDB(), &s, id)
	if err != nil && errors.Is(err, settings.StatusNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update setting: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, s, http.StatusOK)
}

func (sc *SettingsController) DeleteSetting(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = settings.DeleteSettingsById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete setting: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
