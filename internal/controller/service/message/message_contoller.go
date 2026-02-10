package message

import (
	"chickChirick/internal/controller/c_controller"
	"chickChirick/internal/controller/c_http"
	"chickChirick/internal/model/message"
	"encoding/json"
	"errors"
	"net/http"
)

type MessagesController struct {
	Controller c_controller.Controller
}

func (mc *MessagesController) HandleRequest() {
	mc.Controller.ServeMux.HandleFunc("GET /messages", func(w http.ResponseWriter, r *http.Request) {
		mc.GetMessages(w)
	})

	mc.Controller.ServeMux.HandleFunc("POST /message", func(w http.ResponseWriter, r *http.Request) {
		mc.CreateMessage(w, c_http.NewRequest(r))
	})

	mc.Controller.ServeMux.HandleFunc("GET /message/{id}", func(w http.ResponseWriter, r *http.Request) {
		mc.GetMessage(w, c_http.NewRequest(r))
	})

	mc.Controller.ServeMux.HandleFunc("PUT /message/{id}", func(w http.ResponseWriter, r *http.Request) {
		mc.UpdateMessage(w, c_http.NewRequest(r))
	})

	mc.Controller.ServeMux.HandleFunc("DELETE /message/{id}", func(w http.ResponseWriter, r *http.Request) {
		mc.DeleteMessage(w, c_http.NewRequest(r))
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
	id, err := r.HTTPId()
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
	id, err := r.HTTPId()
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
	id, err := r.HTTPId()
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
