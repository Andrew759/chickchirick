package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

type MessagesController struct {
	Controller abstraction.Controller
}

func (mc *MessagesController) HandleRequest() {
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			mc.GetMessages(w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			mc.GetMessage(w, r)
		case http.MethodPost:
			mc.CreateMessage(w, r)
		case http.MethodPut:
			mc.UpdateMessage(w, r)
		case http.MethodDelete:
			mc.DeleteMessage(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (mc *MessagesController) GetMessages(w http.ResponseWriter) {
	messages, err := message.GetMessages(mc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (mc *MessagesController) GetMessage(w http.ResponseWriter, r *http.Request) {
	id := mc.Controller.GETId(w, r)
	m, err := message.GetMessageById(mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Message not found: "+err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (mc *MessagesController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m message.Message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := message.CreateMessage(mc.Controller.Dependencies.DBDecorator.GDB(), &m); err != nil {
		http.Error(w, "Failed to create message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (mc *MessagesController) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m message.Message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := message.UpdateMessage(mc.Controller.Dependencies.DBDecorator.GDB(), &m)
	if err != nil {
		http.Error(w, "Failed to update message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (mc *MessagesController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	id := mc.Controller.GETId(w, r)
	err := message.DeleteMessageById(mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		http.Error(w, "Failed to delete message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
