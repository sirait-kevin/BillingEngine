package usecases

import (
	"context"
	"database/sql"

	"github.com/sirait-kevin/BillingEngine/domain/entities"
	"github.com/sirait-kevin/BillingEngine/domain/interfaces"
)

type DBRepository interface {
	CreateLoan(ctx context.Context, tx *sql.Tx, loan entities.Loan) (int64, error)
	SelectLoanByReferenceId(ctx context.Context, referenceID string) (*entities.Loan, error)
	SelectLoanByUserId(ctx context.Context, userId int64) (*[]entities.Loan, error)
	CreateRepayment(ctx context.Context, tx *sql.Tx, repayment entities.Repayment) (int64, error)
	SelectRepaymentByReferenceId(ctx context.Context, referenceID string) (*entities.Repayment, error)
	SelectRepaymentByLoanId(ctx context.Context, loanIds int64) (*[]entities.Repayment, error)
	SelectTotalRepaymentAmountByLoanId(ctx context.Context, loanId int64) (int64, error)
	SelectRepaymentCountByLoanId(ctx context.Context, loanId int64) (int, error)
	UpdateLoanStatusByReferenceId(ctx context.Context, tx *sql.Tx, referenceId string, status entities.LoanStatus) error

	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type BillingUseCase struct {
	DBRepo DBRepository
	Clock  interfaces.Clock
}
