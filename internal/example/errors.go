package example

import (
	"errors"
	"net/http"
)

var ErrorNotFound error = errors.New("requested item not found")
var ErrorBadInput error = errors.New("input looks wonky")

func getHttpStatusCode(err error) int {
	switch true {
	case errors.Is(err, ErrorBadInput):
		return http.StatusBadRequest
	case errors.Is(err, ErrorNotFound):
		return http.StatusNotFound
	}
	return http.StatusOK
}
