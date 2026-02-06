package c_http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type requestOptions struct {
	requestPrefix string
}

type RequestOption func(options *requestOptions)

// SetRequestPrefix deprecated
func SetRequestPrefix(requestPrefix string) RequestOption {
	return func(rOptions *requestOptions) {
		rOptions.requestPrefix = requestPrefix
	}
}

type Request struct {
	*http.Request
	requestOptions
}

func NewRequest(r *http.Request, opts ...RequestOption) *Request {
	var rOptions requestOptions
	for _, opt := range opts {
		opt(&rOptions)
	}

	return &Request{
		Request:        r,
		requestOptions: rOptions,
	}
}

// HttpId TODO: удалить
// deprecated - достаёт id из URL. Метод предполагает, что ID передается в согласовании с правилами REST API в конце
func (r *Request) HttpId() (int, error) {
	if r.requestPrefix == "" {
		return 0, fmt.Errorf("request prefix not set during request initialization")
	}

	path := strings.TrimPrefix(r.URL.Path, r.requestPrefix)
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		return 0, fmt.Errorf("missing id in URL")
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		return 0, fmt.Errorf("invalid URL id")
	}

	return id, nil
}

func (r *Request) HTTPId() (id int, err error) {
	id, err = strconv.Atoi(r.PathValue("id"))
	if err != nil {
		err = fmt.Errorf("invalid URL")
	}
	return
}
