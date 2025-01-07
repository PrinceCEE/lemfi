package main

import "time"

type Repository struct{}

var usersData = []*User{}
var walletsData = []*UserWallet{}
var transactionsData = []*Transaction{}

func (r *Repository) GetUserByUsername(username string) *User {
	for _, user := range usersData {
		if username == user.Username {
			return user
		}
	}

	return nil
}

func (r *Repository) InsertUser(data *User) {
	usersData = append(usersData, data)
}

func (r *Repository) GetUsers() []*User {
	return usersData
}

func (r *Repository) CreateUserWallet(username string) *UserWallet {
	wallet := &UserWallet{Username: username, Balance: 100000}
	walletsData = append(walletsData, wallet)
	return wallet
}

func (r *Repository) GetUserWallet(username string) *UserWallet {
	for _, wallet := range walletsData {
		if wallet.Username == username {
			return wallet
		}
	}

	return nil
}

func (r *Repository) CreateTransaction(senderUsernamae, receiverUsername string, amount int64) *Transaction {
	t := &Transaction{
		SenderUsername:   senderUsernamae,
		ReceiverUsername: receiverUsername,
		Amount:           amount,
		CreatedAt:        time.Now(),
	}

	transactionsData = append(transactionsData, t)
	return t
}

func (r *Repository) GetUserTransactions(username string) []*Transaction {
	transactions := []*Transaction{}
	for _, t := range transactionsData {
		if t.ReceiverUsername == username || t.SenderUsername == username {
			transactions = append(transactions, t)
		}
	}
	return transactions
}
