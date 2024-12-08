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


## Features

- **Clean Architecture**: Promotes separation of concerns and a decoupled codebase.
- **Middleware**: For logging, error handling, and signature verification to ensure secure API access.
- **Docker**: Containerizes the application for easy setup and deployment.

## Setup and Running the Application

### Prerequisites

- Docker and Docker Compose installed.
- go 1.20 installed

### Steps

1. **Clone the Repository**:
   ```sh
   git clone https://github.com/sirait-kevin/BillingEngine.git
   cd BillingEngine
    ```
2. **Setup**
    ```sh
    run setup.sh
    ```
3. **To Run**
    See the makefile