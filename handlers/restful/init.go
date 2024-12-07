package restful

import (
	"context"

	"github.com/sirait-kevin/BillingEngine/entities"
)

type BillingUsecase interface {
	CreateLoan(ctx context.Context, loanRequest entities.LoanRequest) (int64, error)
	GetLoanHistoryByReferenceID(ctx context.Context, referenceId string) (*entities.LoanHistory, error)
	GetOutStandingAmountByReferenceID(ctx context.Context, referenceId string) (*entities.OutStanding, error)
	GetUserStatusIsDelinquent(ctx context.Context, userId int64) (bool, error)
	GetRepaymentInquiryByLoanReferenceId(ctx context.Context, referenceId string) (*entities.RepaymentInquiry, error)
	MakePayment(ctx context.Context, repaymentRequest entities.RepaymentRequest) (int64, error)
}

type BillingHandler struct {
	BillingUC BillingUsecase
}
