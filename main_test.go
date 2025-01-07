package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	server *httptest.Server
}

func (s *APITestSuite) SetupSuite() {
	repo := &Repository{}
	services := &Services{repo: repo}
	handlers := Handlers{services: services}

	r := chi.NewRouter()
	r.Post("/users", handlers.CreateUser)
	r.Get("/users", handlers.GetUsers)
	r.Post("/transactions", handlers.CreateTransaction)

	s.server = httptest.NewServer(r)
	goroutineCount := 10      // 5 goroutines
	queueInterval := int64(1) // 1s

	for range goroutineCount {
		go HandleVerifyUser(queueInterval)
		go HandleCreateTransaction(queueInterval)
	}
}

func (s *APITestSuite) TearDownSuite() {
	s.server.Close()
}

func (s *APITestSuite) TestAPIEndpoints() {
	client := s.server.Client()
	baseUrl := s.server.URL
	contentType := "application/json"

	users := []*User{
		{
			Username:   "user1",
			IsVerified: false,
		},
		{
			Username:   "user2",
			IsVerified: false,
		},
		{
			Username:   "user3",
			IsVerified: false,
		},
		{
			Username:   "user4",
			IsVerified: false,
		},
		{
			Username:   "user5",
			IsVerified: false,
		},
	}

	s.Run("To create users successfully", func() {
		for _, user := range users {
			createUserDto := CreateUserDto{
				Username: user.Username,
			}
			jsonData, err := json.Marshal(createUserDto)
			s.NoError(err)

			resp, err := client.Post(baseUrl+"/users", contentType, bytes.NewBuffer(jsonData))
			s.NoError(err)
			s.Equal(http.StatusOK, resp.StatusCode)

			var data Response[*User]

			err = ReadJSON(resp.Body, &data)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(true, data.Success)
			s.Equal("User created successfully", data.Message)
			s.Equal(strings.ToLower(user.Username), data.Data.Username)
		}
	})

	s.Run("To fetch users successfully", func() {
		resp, err := client.Get(baseUrl + "/users")
		s.NoError(err)
		s.Equal(http.StatusOK, resp.StatusCode)

		var data Response[[]*User]
		err = ReadJSON(resp.Body, &data)
		s.NoError(err)
		defer resp.Body.Close()

		s.Equal(true, data.Success)
		s.Equal("Fetched users successfully", data.Message)
		s.Equal(5, len(data.Data))

		for _, user := range data.Data {
			s.Equal(true, user.IsVerified)
			s.Equal(0, len(user.Transactions))
			s.Equal(int64(100000), user.Wallet.Balance)
		}
	})

	s.Run("To create transactions", func() {
		sender := users[0]
		receiver := users[1]

		createTransactionDto := CreateTransactionDto{
			ReceiverUsername: receiver.Username,
			SenderUsername:   sender.Username,
			Amount:           5000,
		}

		jsonData, err := json.Marshal(createTransactionDto)
		s.NoError(err)

		resp, err := client.Post(baseUrl+"/transactions", contentType, bytes.NewBuffer(jsonData))
		s.NoError(err)
		defer resp.Body.Close()

		var data Response[*Transaction]
		err = ReadJSON(resp.Body, &data)
		s.NoError(err)

		s.Equal(true, data.Success)
		s.Equal("Created transactions successfully", data.Message)
		s.Equal(int64(5000), data.Data.Amount)
		s.Equal(sender.Username, data.Data.SenderUsername)
		s.Equal(receiver.Username, data.Data.ReceiverUsername)

		s.Run("To fetch users successfully after transaction", func() {
			resp, err := client.Get(baseUrl + "/users")
			s.NoError(err)
			s.Equal(http.StatusOK, resp.StatusCode)

			var data Response[[]*User]
			err = ReadJSON(resp.Body, &data)
			s.NoError(err)
			defer resp.Body.Close()

			s.Equal(true, data.Success)
			s.Equal("Fetched users successfully", data.Message)
			s.Equal(5, len(data.Data))

			for _, user := range data.Data {
				if user.Username == sender.Username {
					s.Equal(1, len(user.Transactions))
					s.Equal(int64(95000), user.Wallet.Balance)
				}
				if user.Username == receiver.Username {
					s.Equal(1, len(user.Transactions))
					s.Equal(int64(105000), user.Wallet.Balance)
				}
			}
		})
	})
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
