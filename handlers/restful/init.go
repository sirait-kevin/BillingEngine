package restful

import (
	"context"

	entities2 "github.com/sirait-kevin/BillingEngine/domain/entities"
)

type BillingUsecase interface {
	CreateLoan(ctx context.Context, loanRequest entities2.LoanRequest) (int64, error)
	GetLoanHistoryByReferenceID(ctx context.Context, referenceId string) (*entities2.LoanHistory, error)
	GetOutStandingAmountByReferenceID(ctx context.Context, referenceId string) (*entities2.OutStanding, error)
	GetUserStatusIsDelinquent(ctx context.Context, userId int64) (bool, error)
	GetRepaymentInquiryByLoanReferenceId(ctx context.Context, referenceId string) (*entities2.RepaymentInquiry, error)
	MakePayment(ctx context.Context, repaymentRequest entities2.RepaymentRequest) (int64, error)
}

type BillingHandler struct {
	BillingUC BillingUsecase
}
