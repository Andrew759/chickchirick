package http_transaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Response struct {
	*ResultContainer
	ErrorContainer []*ErrorBody `json:"errors,omitempty"`
}

type ResultContainer struct {
	Result interface{}   `json:"result,omitempty"`
	buf    *bytes.Reader // Внутренний буфер для чтения
}

type ErrorBody struct {
	Message string `json:"message"`
}

func NewResponse() *Response {
	return &Response{}
}

// SendSuccess перезаписывает ResultContainer и отправляет ответ
func (r *Response) SendSuccess(w http.ResponseWriter, code int, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	r.ResultContainer = &ResultContainer{
		Result: result,
	}

	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		r.SendError(w, http.StatusInternalServerError, err.Error())
	}
}

// SendError добавляет одну ошибку в структуру ошибок и отправляет ответ
func (r *Response) SendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errBody := ErrorBody{
		Message: message,
	}
	r.ErrorContainer = append(r.ErrorContainer, &errBody)

	_ = json.NewEncoder(w).Encode(r)
}

// AddErrorToErrorContainer добавление ошибки в контейнер ошибок
func (r *Response) AddErrorToErrorContainer(err error) {
	r.ErrorContainer = append(r.ErrorContainer, &ErrorBody{err.Error()})
}

func (r *Response) AddErrorsToErrorContainer(errors []error) {
	for _, err := range errors {
		r.AddErrorToErrorContainer(err)
	}
}

// Send отправка данных, записанных в структуру Response
func (r *Response) Send(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(r)
}

// FirstError Возвращает первую ошибку из ErrorContainer
func (r *Response) FirstError() error {
	if len(r.ErrorContainer) > 0 {
		return errors.New(r.ErrorContainer[0].Message)
	}
	return nil
}

func (rc *ResultContainer) Read(p []byte) (n int, err error) {
	if rc.buf == nil {
		data, err := json.Marshal(rc.Result)
		if err != nil {
			return 0, err
		}
		rc.buf = bytes.NewReader(data)
	}
	return rc.buf.Read(p)
}
