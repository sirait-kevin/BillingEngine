package usecases

import (
	"context"
	"math"
	"net/http"
	"strings"
	"sync"

	"github.com/sirait-kevin/BillingEngine/entities"
	"github.com/sirait-kevin/BillingEngine/pkg/errs"
	"github.com/sirait-kevin/BillingEngine/pkg/helper"
)

func (u *BillingUseCase) CreateLoan(ctx context.Context, loanRequest entities.Loan) (int64, error) {
	var errMessage []string

	if loanRequest.ReferenceId == "" {
		errMessage = append(errMessage, "Reference Id is required")
	}
	if !IsUserValid(loanRequest.UserId) {
		errMessage = append(errMessage, "UserId is invalid")
	}
	if loanRequest.Amount < 1 {
		errMessage = append(errMessage, "Amount is required")
	}
	if loanRequest.RatePercentage < 0 {
		errMessage = append(errMessage, "Rate percentage is required")
	}
	if !loanRequest.RepaymentSchedule.IsValid() {
		errMessage = append(errMessage, "Repayment schedule is invalid")
	}
	if loanRequest.Tenor < 1 {
		errMessage = append(errMessage, "Tenor is required")
	}

	if errMessage != nil || len(errMessage) != 0 {
		return 0, errs.NewWithMessage(http.StatusBadRequest, strings.Join(errMessage, ","))
	}

	repaymentAmount := int64(math.Round(float64(loanRequest.Amount / int64(loanRequest.Tenor))))

	loanId, err := u.DBRepo.CreateLoan(ctx, entities.Loan{
		ReferenceId:       loanRequest.ReferenceId,
		UserId:            loanRequest.UserId,
		Amount:            loanRequest.Amount,
		RatePercentage:    loanRequest.RatePercentage,
		Status:            entities.LoanStatusActive,
		RepaymentSchedule: loanRequest.RepaymentSchedule,
		Tenor:             loanRequest.Tenor,
		RepaymentAmount:   repaymentAmount,
	})
	if err != nil {
		return 0, err
	}
	return loanId, nil

}

func (u *BillingUseCase) GetLoanHistoryByReferenceID(ctx context.Context, referenceId string) (*entities.LoanHistory, error) {

	if referenceId == "" {
		return nil, errs.NewWithMessage(http.StatusBadRequest, "reference id can not be empty")
	}

	loan, err := u.DBRepo.SelectLoanByReferenceId(ctx, referenceId)
	if err != nil {
		return nil, err
	}

	repayments, err := u.DBRepo.SelectRepaymentByLoanId(ctx, loan.Id)
	if err != nil {
		if errs.GetHTTPCode(err) != http.StatusNotFound {
			return nil, err
		}
	}

	return &entities.LoanHistory{
		Loan:       *loan,
		Repayments: *repayments,
	}, nil
}

func (u *BillingUseCase) GetOutStandingAmountByReferenceID(ctx context.Context, referenceId string) (*entities.OutStanding, error) {
	if referenceId == "" {
		return nil, errs.NewWithMessage(http.StatusBadRequest, "reference id can not be empty")
	}
	loan, err := u.DBRepo.SelectLoanByReferenceId(ctx, referenceId)
	if err != nil {
		return nil, err
	}
	totalRepayments, err := u.DBRepo.SelectTotalRepaymentAmountByLoanId(ctx, loan.Id)
	if err != nil {
		if errs.GetHTTPCode(err) != http.StatusNotFound {
			return nil, err
		}
	}

	totalLoan := loan.RepaymentAmount * int64(loan.Tenor)

	return &entities.OutStanding{
		LoanId:            loan.Id,
		LoanReferenceId:   loan.ReferenceId,
		OutstandingAmount: totalLoan - totalRepayments,
	}, nil
}

func (u *BillingUseCase) GetUserStatusIsDelinquent(ctx context.Context, userId int64) (bool, error) {
	if IsUserValid(userId) {
		return false, errs.NewWithMessage(http.StatusBadRequest, "user id is invalid")
	}

	loans, err := u.DBRepo.SelectLoanByUserId(ctx, userId)
	if err != nil {
		return false, err
	}
	repaymentCounts := map[int64]int{}
	errWg := make([]error, len(*loans))

	var wg sync.WaitGroup
	for i, loan := range *loans {
		wg.Add(1)
		go func(i int, loan entities.Loan) {
			defer wg.Done()
			repaymentCounts[loan.Id], errWg[i] = u.DBRepo.SelectRepaymentCountByLoanId(ctx, loan.Id)
			if errWg[i] != nil {
				if errs.GetHTTPCode(errWg[i]) != http.StatusNotFound {
					errWg[i] = nil
				}
			}
		}(i, loan)
	}
	wg.Wait()

	for _, err = range errWg {
		if err != nil {
			return false, err
		}
	}

	for _, loan := range *loans {
		if loan.Tenor > repaymentCounts[loan.Id] {
			if repaymentCounts[loan.Id] < helper.MonthsBetween(loan.CreatedAt, u.Clock.Now()) {
				return true, nil
			}
		}
	}

	return false, nil
}

func (u *BillingUseCase) GetRepaymentInquiryByLoanReferenceId(ctx context.Context, referenceId string) (*entities.RepaymentInquiry, error) {
	if referenceId == "" {
		return nil, errs.NewWithMessage(http.StatusBadRequest, "reference id can not be empty")
	}

	loan, err := u.DBRepo.SelectLoanByReferenceId(ctx, referenceId)
	if err != nil {
		return nil, err
	}

	repaymentCount, err := u.DBRepo.SelectRepaymentCountByLoanId(ctx, loan.Id)
	if err != nil {
		if errs.GetHTTPCode(err) != http.StatusNotFound {
			return nil, err
		}
	}

	missRepaymentCount := helper.MonthsBetween(loan.CreatedAt, u.Clock.Now()) + 1 - repaymentCount
	if missRepaymentCount == 0 {
		return &entities.RepaymentInquiry{
			LoanId:          loan.Id,
			LoanReferenceId: loan.ReferenceId,
			RepaymentNeeded: nil,
		}, nil
	}

	needRepayments := make([]entities.RepaymentNeeded, missRepaymentCount)
	for i := 0; i < len(needRepayments); i++ {

		dueDate := loan.CreatedAt.AddDate(0, i, 0)
		isLate := u.Clock.Now().After(dueDate)

		needRepayments[i] = entities.RepaymentNeeded{
			Amount:  loan.RepaymentAmount,
			DueDate: dueDate,
			IsLate:  isLate,
		}
	}

	return &entities.RepaymentInquiry{
		LoanId:          loan.Id,
		LoanReferenceId: referenceId,
		RepaymentNeeded: needRepayments,
	}, nil
}

func (u *BillingUseCase) MakePayment(ctx context.Context, repaymentRequest entities.Repayment) error {
	var errMessage []string

	if repaymentRequest.LoanId < 1 {
		errMessage = append(errMessage, "loan id is invalid")
	}
	if repaymentRequest.Amount < 1 {
		errMessage = append(errMessage, "amount is invalid")
	}
	if repaymentRequest.ReferenceId == "" {
		errMessage = append(errMessage, "reference id can not be empty")
	}
	if errMessage != nil || len(errMessage) != 0 {
		return errs.NewWithMessage(http.StatusBadRequest, strings.Join(errMessage, "; "))
	}

	err := u.MakePayment(ctx, repaymentRequest)
	if err != nil {
		return err
	}
	return nil
}

func IsUserValid(userId int64) bool {
	if userId < 1 {
		return false
	}

	//TODO: this function has to hit an account service API.
	return false
}
