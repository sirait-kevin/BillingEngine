# BillingEngine
<pre>
BillingEngine
├── Makefile
├── README.md
├── docker-compose.yml
├── domain
│   ├── entities
│   │   ├── request.go
│   │   ├── response.go
│   │   └── transaction.go
│   └── interfaces
│       ├── clock.go
│       └── sql.go
├── go.mod
├── go.sum
├── handlers
│   ├── middleware
│   │   └── middleware.go
│   ├── mq
│   │   └── init.go
│   └── restful
│       ├── billing.go
│       ├── billing_test.go
│       └── init.go
├── main
│   └── main.go
├── mocks
│   ├── domain
│   │   ├── atomic_transaction.go
│   │   └── clock.go
│   ├── handler
│   │   └── BillingUsecase.go
│   └── usecases
│       └── DBRepository.go
├── pkg
│   ├── errs
│   │   └── errors.go
│   ├── helper
│   │   ├── clock.go
│   │   ├── encryptions.go
│   │   ├── parser.go
│   │   └── response.go
│   └── logger
│       └── logger.go
├── repositories
│   ├── data_type.go
│   ├── init.go
│   └── transaction.go
├── setup.sh
├── sqlfiles
│   └── init.sql
├── usecases
│   ├── init.go
│   ├── transactions.go
│   └── transactions_test.go
└── vendor
</pre>

