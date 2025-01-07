package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func WriteResponse(w http.ResponseWriter, msg string, data any) error {
	resp := Response[any]{
		Success: true,
		Message: msg,
		Data:    data,
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	return nil
}

func WriteErrorResponse(w http.ResponseWriter, err error, statusCode int) {
	resp := Response[any]{
		Success: false,
		Message: err.Error(),
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("unexpected error occured: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

func ReadJSON(body io.Reader, data any) error {
	d := json.NewDecoder(body)
	return d.Decode(data)
}
