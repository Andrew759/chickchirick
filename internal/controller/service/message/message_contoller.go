package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/model/message"
	"encoding/json"
	"net/http"
)

const MessageResource = "/message/"

type MessagesController struct {
	Controller abstraction.Controller
}

func (mc *MessagesController) HandleRequest() {
	mc.Controller.ServeMux.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			mc.GetMessages(w)
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mc.Controller.ServeMux.HandleFunc(MessageResource, func(w http.ResponseWriter, r *http.Request) {
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
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func (mc *MessagesController) GetMessages(w http.ResponseWriter) {
	messages, err := message.GetMessages(mc.Controller.Dependencies.DBDecorator.GDB())
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, messages, http.StatusOK)
}

func (mc *MessagesController) GetMessage(w http.ResponseWriter, r *http.Request) {
	id := mc.Controller.HttpId(w, r, MessageResource)
	m, err := message.GetMessageById(mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Message not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, m, http.StatusOK)
}

func (mc *MessagesController) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var m message.Message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := message.CreateMessage(mc.Controller.Dependencies.DBDecorator.GDB(), &m); err != nil {
		c_http.NewResponse().SendError(w, "Failed to create message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, m, http.StatusCreated)
}

func (mc *MessagesController) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	var m message.Message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err := message.UpdateMessage(mc.Controller.Dependencies.DBDecorator.GDB(), &m)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, m, http.StatusOK)
}

func (mc *MessagesController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	id := mc.Controller.HttpId(w, r, MessageResource)
	err := message.DeleteMessageById(mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
