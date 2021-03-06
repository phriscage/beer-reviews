package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ResponseCore struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

type ResponseData struct {
	ResponseCore
	Data interface{} `json:"data,omitempty"`
}

type ResponseError struct {
	ResponseCore
	Errors []string `json:"errors,omitempty"`
}

// Custom Response Error Handlers with errors array
func ResponseErrorHandler(w http.ResponseWriter, r *http.Request, code int, errors []string) {
	if len(errors) != 0 {
		log.Warn(errors)
	}
	responseData := &ResponseError{ResponseCore{code, http.StatusText(code)}, errors}
	body, err := json.MarshalIndent(responseData, "", "    ")
	if err != nil {
		log.Fatal(err)
		return
	}
	// Always set Headers before Writing them
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// add new line to response
	w.Write(append([]byte(body), []byte("\n")...))
}

// Custom Response Handlers with data array
func ResponseHandler(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	responseData := &ResponseData{ResponseCore{code, http.StatusText(code)}, data}
	body, err := json.MarshalIndent(responseData, "", "    ")
	if err != nil {
		log.Fatal(err)
		return
	}
	// Always set Headers before Writing them
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// add new line to response
	w.Write(append([]byte(body), []byte("\n")...))
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	ResponseErrorHandler(w, r, http.StatusMethodNotAllowed, nil)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	ResponseHandler(w, r, http.StatusNotFound, nil)
}
