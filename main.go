package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	var queueInterval int64
	var goroutineCount int64
	var port int64

	flag.Int64Var(&queueInterval, "queueInterval", 5, "The verification interval")
	flag.Int64Var(&goroutineCount, "goroutineCount", 5, "Number of goroutines")
	flag.Int64Var(&port, "port", 3000, "Port to run server")
	flag.Parse()

	repo := &Repository{}
	services := &Services{repo: repo}
	handlers := Handlers{services: services}

	r := chi.NewRouter()
	r.Post("/users", handlers.CreateUser)
	r.Get("/users", handlers.GetUsers)
	r.Post("/transactions", handlers.CreateTransaction)

	for range goroutineCount {
		go HandleVerifyUser(queueInterval)
		go HandleCreateTransaction(queueInterval)
	}

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	fmt.Printf("Starting server on port :%d\n", port)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
