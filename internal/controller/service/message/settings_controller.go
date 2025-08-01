package message

import (
	"chickChirick/internal/controller/abstraction"
	settings "chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

type SettingsController struct {
	Controller abstraction.Controller
}

func (sc *SettingsController) HandleRequest(mux *http.ServeMux) {
	mux.HandleFunc("/message/settings", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sc.GetSettings(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/message/setting", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sc.GetSetting(w, r)
		case http.MethodPost:
			sc.CreateSetting(w, r)
		case http.MethodPut:
			sc.UpdateSetting(w, r)
		case http.MethodDelete:
			sc.DeleteSetting(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (sc *SettingsController) GetSettings(w http.ResponseWriter) {
	settingList, err := settings.GetSettings(sc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(settingList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sc *SettingsController) GetSetting(w http.ResponseWriter, r *http.Request) {
	id := sc.Controller.GETId(w, r)
	s, err := settings.GetSettingsById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Setting not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sc *SettingsController) CreateSetting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var s settings.Settings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := settings.CreateSettings(sc.Controller.Dependencies.DBDecorator.GDB(), &s); err != nil {
		http.Error(w, "Failed to create setting: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sc *SettingsController) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var s settings.Settings
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := settings.UpdateSettings(sc.Controller.Dependencies.DBDecorator.GDB(), &s)
	if err != nil {
		http.Error(w, "Failed to update setting: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(s); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sc *SettingsController) DeleteSetting(w http.ResponseWriter, r *http.Request) {
	id := sc.Controller.GETId(w, r)
	err := settings.DeleteSettingsById(sc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete setting: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
