package http_transaction

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Response struct {
	*ResultContainer
	ErrorContainer *ErrorBody `json:"error,omitempty"`
}

type ResultContainer struct {
	Result interface{}   `json:"result,omitempty"`
	buf    *bytes.Reader // Внутренний буфер для чтения
}

type ErrorBody struct {
	Message string `json:"message"`
}

func NewResponse() Response {
	return Response{}
}

// Success отправка ответа
func (r Response) Success(w http.ResponseWriter, code int, result interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := Response{
		ResultContainer: &ResultContainer{
			Result: result,
		},
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		r.Error(w, http.StatusInternalServerError, err.Error())
	}
}

// Error отправка ошибки
func (r Response) Error(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := Response{
		ErrorContainer: &ErrorBody{
			Message: message,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
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
