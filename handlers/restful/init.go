package restful

import (
	"context"

	"github.com/sirait-kevin/BillingEngine/domain/entities"
)

type BillingUsecase interface {
	CreateLoan(ctx context.Context, loanRequest entities.LoanRequest) (int64, error)
	GetPaymentHistoryByReferenceID(ctx context.Context, referenceId string) (*entities.LoanHistory, error)
	GetOutStandingAmountByReferenceID(ctx context.Context, referenceId string) (*entities.OutStanding, error)
	GetUserStatusIsDelinquent(ctx context.Context, userId int64) (bool, error)
	GetRepaymentInquiryByLoanReferenceId(ctx context.Context, referenceId string) (*entities.RepaymentInquiry, error)
	MakePayment(ctx context.Context, repaymentRequest entities.RepaymentRequest) (int64, error)
	GetLoanListByUserId(ctx context.Context, userId int64) (*[]entities.Loan, error)
}

type BillingHandler struct {
	BillingUC BillingUsecase
}
