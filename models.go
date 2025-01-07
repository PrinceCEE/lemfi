package main

import "time"

type User struct {
	Username     string         `json:"username"`
	IsVerified   bool           `json:"is_verified"`
	Wallet       *UserWallet    `json:"wallet"`
	Transactions []*Transaction `json:"transactions"`
}

type UserWallet struct {
	Username string `json:"username"`
	Balance  int64  `json:"balance"`
}

type Transaction struct {
	ReceiverUsername string    `json:"receiver_username"`
	SenderUsername   string    `json:"sender_username"`
	Amount           int64     `json:"amount"`
	CreatedAt        time.Time `json:"created_at"`
}
