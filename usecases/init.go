package usecases

import (
	"context"

	"github.com/sirait-kevin/BillingEngine/entities"
)

type DBRepository interface {
	GetByID(ctx context.Context, id int64) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) (int64, error)
	Update(ctx context.Context, user *entities.User) error
	CreateLoan(ctx context.Context, loan *entities.Loan) (int64, error)
	SelectLoanByReferenceId(ctx context.Context, referenceID string) (*entities.Loan, error)
	SelectLoanByUserId(ctx context.Context, userId int64) (*[]entities.Loan, error)
	CreateRepayment(ctx context.Context, repayment *entities.Repayment) (int64, error)
	SelectRepaymentByReferenceId(ctx context.Context, referenceID string) (*entities.Repayment, error)
	SelectRepaymentByLoanId(ctx context.Context, referenceID string) (*[]entities.Repayment, error)
}

type UserUseCase struct {
	DBRepo DBRepository
}
