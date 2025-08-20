package message

import (
	"chickChirick/internal/controller/abstraction"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/model/message"
	"encoding/json"
	"errors"
	"net/http"
)

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

	mc.Controller.ServeMux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			mc.CreateMessage(w, c_http.NewRequest(r))
		default:
			c_http.NewResponse().SendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mc.Controller.ServeMux.HandleFunc("/message/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			mc.GetMessage(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/message/")))
		case http.MethodPut:
			mc.UpdateMessage(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/message/")))
		case http.MethodDelete:
			mc.DeleteMessage(w, c_http.NewRequest(r, c_http.SetRequestPrefix("/message/")))
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

func (mc *MessagesController) GetMessage(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := message.GetMessageById(mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil {
		c_http.NewResponse().SendError(w, "Message not found: "+err.Error(), http.StatusNotFound)
		return
	}

	c_http.NewResponse().SendSuccess(w, m, http.StatusOK)
}

func (mc *MessagesController) CreateMessage(w http.ResponseWriter, r *c_http.Request) {
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

func (mc *MessagesController) UpdateMessage(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var m message.Message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		c_http.NewResponse().SendError(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = message.UpdateMessageById(mc.Controller.Dependencies.DBDecorator.GDB(), &m, id)
	if err != nil && errors.Is(err, message.MessageNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to update message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	c_http.NewResponse().SendSuccess(w, m, http.StatusOK)
}

func (mc *MessagesController) DeleteMessage(w http.ResponseWriter, r *c_http.Request) {
	id, err := r.HttpId()
	if err != nil {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = message.DeleteMessageById(mc.Controller.Dependencies.DBDecorator.GDB(), id)
	if err != nil && errors.Is(err, message.MessageNotFoundErr) {
		c_http.NewResponse().SendError(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		c_http.NewResponse().SendError(w, "Failed to delete message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
