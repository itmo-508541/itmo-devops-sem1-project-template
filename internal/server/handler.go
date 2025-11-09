package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	contentTypeHeader = "Content-Type"
	jsonContentType   = "application/json; charset=UTF-8"
)

type (
	Handler[TIn, TOut any] func(TIn, *TOut) error

	errorResponseDTO struct {
		Error errorDTO `json:"error"`
	}

	errorDTO struct {
		Message string `json:"message"`
	}
)

func (h Handler[TIn, TOut]) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Println(panicErr)
		}
	}()

	writer.Header().Set(contentTypeHeader, jsonContentType)

	var in TIn
	if request.Body != nil && request.Body != http.NoBody {
		if err := json.NewDecoder(request.Body).Decode(&in); err != nil {
			log.Println(fmt.Errorf("json.Decode: %w", err).Error())
			badRequestError(writer)

			return
		}
	}

	var out TOut
	if err := h(in, &out); err != nil {
		log.Println(fmt.Errorf("h: %w", err).Error())
		internalServerError(writer)

		return
	}

	if err := json.NewEncoder(writer).Encode(out); err != nil {
		log.Println(fmt.Errorf("json.Encode: %w", err).Error())
	}
}

func baseError(writer http.ResponseWriter, message string, code int) {
	response := errorResponseDTO{
		Error: errorDTO{
			Message: message,
		},
	}

	writer.WriteHeader(code)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Println(fmt.Errorf("json.Encode: %w", err).Error())
	}
}

func badRequestError(writer http.ResponseWriter) {
	baseError(writer, "Bad Request", http.StatusBadRequest)
}

func internalServerError(writer http.ResponseWriter) {
	baseError(writer, "Internal Server Error", http.StatusInternalServerError)
}
