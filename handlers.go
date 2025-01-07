package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Handlers struct {
	services *Services
}

var v = validator.New()

type CreateUserDto struct {
	Username string `json:"username" validator:"alpha"`
}

type CreateTransactionDto struct {
	ReceiverUsername string `json:"receiver_username"`
	SenderUsername   string `json:"sender_username"`
	Amount           int64  `json:"amount"`
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var createUserDto CreateUserDto

	err := ReadJSON(r.Body, &createUserDto)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = v.Struct(createUserDto)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	res, code, err := h.services.CreateUser(&createUserDto)
	if err != nil {
		WriteErrorResponse(w, err, code)
		return
	}

	WriteResponse(w, "User created successfully", res)
}

func (h *Handlers) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, code, err := h.services.GetUsers()
	if err != nil {
		WriteErrorResponse(w, err, code)
	}

	WriteResponse(w, "Fetched users successfully", users)
}

func (h *Handlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var createTransactionDto CreateTransactionDto
	err := ReadJSON(r.Body, &createTransactionDto)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = v.Struct(createTransactionDto)
	if err != nil {
		WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	t, code, err := h.services.CreateTransaction(&createTransactionDto)
	if err != nil {
		WriteErrorResponse(w, err, code)
		return
	}

	WriteResponse(w, "Created transactions successfully", t)
}
