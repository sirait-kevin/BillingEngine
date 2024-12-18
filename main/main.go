package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nsqio/go-nsq"

	"github.com/sirait-kevin/BillingEngine/handlers/middleware"
	"github.com/sirait-kevin/BillingEngine/handlers/mq"
	"github.com/sirait-kevin/BillingEngine/handlers/restful"
	"github.com/sirait-kevin/BillingEngine/pkg/helper"
	"github.com/sirait-kevin/BillingEngine/pkg/logger"
	"github.com/sirait-kevin/BillingEngine/repositories"
	"github.com/sirait-kevin/BillingEngine/usecases"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize logger
	logger.InitLogger(true)
	logger.Info("Starting BillingEngine...")

	db, err := sqlx.Connect("mysql", "BillingEngine:rootpassword@tcp(localhost:3306)/BillingEngine?parseTime=true")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	dbRepository := &repositories.DBRepository{DB: db}
	billingUsecase := &usecases.BillingUseCase{
		DBRepo: dbRepository,
		Clock:  helper.RealClock{},
	}
	billingHandler := &restful.BillingHandler{BillingUC: billingUsecase}

	router := mux.NewRouter()

	router.Use(middleware.VerifySignatureMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.ErrorHandlingMiddleware)

	router.HandleFunc("/create/loan", billingHandler.CreateLoan).Methods(http.MethodPost)
	router.HandleFunc("/make/payment", billingHandler.MakePayment).Methods(http.MethodPost)

	router.HandleFunc("/payment/history", billingHandler.GetPaymentHistory).Methods(http.MethodGet)
	router.HandleFunc("/outstanding/amount", billingHandler.GetOutStandingAmount).Methods(http.MethodGet)
	router.HandleFunc("/user/status", billingHandler.GetUserStatus).Methods(http.MethodGet)
	router.HandleFunc("/payment/inquiry", billingHandler.GetPaymentInquiry).Methods(http.MethodGet)
	router.HandleFunc("/loan/history", billingHandler.GetLoanHistory).Methods(http.MethodGet)

	//nsqHandler := &mq.NSQHandler{BillingUseCase: useCase}
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
