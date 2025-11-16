package server

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	contentTypeHeader        = "Content-Type"
	contentDispositionHeader = "Content-Disposition"
	jsonContentType          = "application/json; charset=UTF-8"
	textContentType          = "plain/text; charset=UTF-8"
	zipContentType           = "application/zip"
)

// https://stackoverflow.com/questions/46791169/create-serve-over-http-a-zip-file-without-writing-to-disk
func ZipResponse(
	writer http.ResponseWriter,
	responseName string,
	contents string,
	contentsName string,
	code int,
) {
	var err error

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	f, err := zw.Create(contentsName)
	if err == nil {
		_, err = f.Write([]byte(contents))
	}
	if err == nil {
		err = zw.Close()
	}
	if err != nil {
		log.Printf("ZipResponse: %s\n", err.Error())
		JSONInternalServerError(writer)
		return
	}

	writer.WriteHeader(code)
	writer.Header().Set(contentTypeHeader, zipContentType)
	writer.Header().
		Set(contentDispositionHeader, fmt.Sprintf("attachment; filename=\"%s\"", responseName))
	writer.Write(buf.Bytes())
}

func TextResponse(writer http.ResponseWriter, response string, code int) {
	writer.WriteHeader(code)
	writer.Header().Set(contentTypeHeader, textContentType)
	writer.Write([]byte(response))
}

func JSONResponse(writer http.ResponseWriter, response any, code int) {
	writer.WriteHeader(code)
	writer.Header().Set(contentTypeHeader, jsonContentType)

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Println(fmt.Errorf("json.Encode: %w", err).Error())
	}
}

func JSONBaseError(writer http.ResponseWriter, message string, code int) {
	response := errorResponseDTO{
		Error: errorDTO{
			Message: message,
		},
	}
	JSONResponse(writer, response, code)
}

func JSONBadRequestError(writer http.ResponseWriter) {
	JSONBaseError(writer, "Bad Request", http.StatusBadRequest)
}

func JSONInternalServerError(writer http.ResponseWriter) {
	JSONBaseError(writer, "Internal Server Error", http.StatusInternalServerError)
}
