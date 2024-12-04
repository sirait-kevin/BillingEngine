package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nsqio/go-nsq"

	_ "github.com/go-sql-driver/mysql"

	"github.com/sirait-kevin/BillingEngine/handlers/middleware"
	"github.com/sirait-kevin/BillingEngine/handlers/mq"
	"github.com/sirait-kevin/BillingEngine/handlers/restful"
	"github.com/sirait-kevin/BillingEngine/pkg/logger"
	"github.com/sirait-kevin/BillingEngine/repositories"
	"github.com/sirait-kevin/BillingEngine/usecases"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize logger
	logger.InitLogger(true)
	logger.Info("Starting BillingEngine...")

	db, err := sql.Open("mysql", "BillingEngine:rootpassword@tcp(localhost:3306)/BillingEngine")
	if err != nil {
		log.Fatal("error open sql", err)
	}
	defer db.Close()

	userRepository := &repositories.UserRepository{DB: db}
	userUseCase := &usecases.UserUseCase{UserRepository: userRepository}
	userHandler := &restful.UserHandler{UserUseCase: userUseCase}

	router := mux.NewRouter()

	router.Use(middleware.VerifySignatureMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.ErrorHandlingMiddleware)

	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")

	//nsqHandler := &mq.NSQHandler{UserUseCase: userUseCase}
	//startNSQConsumer(nsqHandler)

	logger.Log.Info("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func startNSQConsumer(handler *mq.NSQHandler) {
	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("user_updates", "channel", config)
	q.AddHandler(nsq.HandlerFunc(handler.HandleMessage))
	err := q.ConnectToNSQLookupd("localhost:4161")
	if err != nil {
		log.Panic("Could not connect to NSQ")
	}
}
