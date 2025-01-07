package main

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

type Services struct {
	repo *Repository
}

var verifyQueue = make(chan string)
var transactionQueue = make(chan *Transaction)

func HandleVerifyUser(sleep int64) {
	repo := &Repository{}

	for {
		username := <-verifyQueue
		user := repo.GetUserByUsername(username)
		if user != nil {
			user.IsVerified = true
			repo.CreateUserWallet(username)
		}

		time.Sleep(time.Second * time.Duration(sleep))
	}
}

func HandleCreateTransaction(sleep int64) {
	repo := &Repository{}
	for {
		t := <-transactionQueue
		senderWallet := repo.GetUserWallet(t.SenderUsername)
		receicerWallet := repo.GetUserWallet(t.ReceiverUsername)
		senderWallet.Balance -= t.Amount
		receicerWallet.Balance += t.Amount

		repo.CreateTransaction(t.SenderUsername, t.ReceiverUsername, t.Amount)
		time.Sleep(time.Second * time.Duration(sleep))
	}
}

func (s *Services) CreateUser(data *CreateUserDto) (*User, int, error) {
	data.Username = strings.ToLower(data.Username)
	if data.Username == "" {
		return nil, http.StatusBadRequest, errors.New("username missing")
	}

	userExists := s.repo.GetUserByUsername(data.Username)
	if userExists != nil {
		return nil, http.StatusBadRequest, errors.New("user already exists")
	}

	user := &User{Username: data.Username, IsVerified: false}
	s.repo.InsertUser(user)

	verifyQueue <- data.Username

	return user, http.StatusOK, nil
}

func (s *Services) GetUsers() ([]*User, int, error) {
	users := s.repo.GetUsers()
	for _, user := range users {
		wallet := s.repo.GetUserWallet(user.Username)
		user.Wallet = wallet
		user.Transactions = s.repo.GetUserTransactions(user.Username)
	}
	return users, http.StatusOK, nil
}

func (s *Services) CreateTransaction(data *CreateTransactionDto) (*Transaction, int, error) {
	data.SenderUsername = strings.ToLower(data.SenderUsername)
	data.ReceiverUsername = strings.ToLower(data.ReceiverUsername)

	sender := s.repo.GetUserByUsername(data.SenderUsername)
	receiver := s.repo.GetUserByUsername(data.ReceiverUsername)
	if sender == nil {
		return nil, http.StatusNotFound, errors.New("sender user not found")
	}
	if receiver == nil {
		return nil, http.StatusNotFound, errors.New("receiver user not found")
	}

	if !sender.IsVerified {
		return nil, http.StatusUnauthorized, errors.New("sender user not verified")
	}
	if !receiver.IsVerified {
		return nil, http.StatusUnauthorized, errors.New("receiver user not verified")
	}

	senderWallet := s.repo.GetUserWallet(data.SenderUsername)
	if senderWallet.Balance-data.Amount < 0 {
		return nil, http.StatusBadRequest, errors.New("insufficient balance")
	}

	t := &Transaction{
		ReceiverUsername: data.ReceiverUsername,
		SenderUsername:   data.SenderUsername,
		Amount:           data.Amount,
	}
	transactionQueue <- t

	return t, http.StatusOK, nil
}
