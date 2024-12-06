package usecases

import (
	"context"
	"math"
	"net/http"
	"strings"
	"sync"

	entities2 "github.com/sirait-kevin/BillingEngine/domain/entities"
	"github.com/sirait-kevin/BillingEngine/pkg/errs"
	"github.com/sirait-kevin/BillingEngine/pkg/helper"
)

func (u *BillingUseCase) CreateLoan(ctx context.Context, loanRequest entities2.LoanRequest) (int64, error) {
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

	loanId, err := u.DBRepo.CreateLoan(ctx, entities2.Loan{
		ReferenceId:       loanRequest.ReferenceId,
		UserId:            loanRequest.UserId,
		Amount:            loanRequest.Amount,
		RatePercentage:    loanRequest.RatePercentage,
		Status:            entities2.LoanStatusActive,
		RepaymentSchedule: loanRequest.RepaymentSchedule,
		Tenor:             loanRequest.Tenor,
		RepaymentAmount:   repaymentAmount,
	})
	if err != nil {
		return 0, err
	}
	return loanId, nil

}

func (u *BillingUseCase) GetLoanHistoryByReferenceID(ctx context.Context, referenceId string) (*entities2.LoanHistory, error) {

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

	return &entities2.LoanHistory{
		Loan:       *loan,
		Repayments: *repayments,
	}, nil
}

func (u *BillingUseCase) GetOutStandingAmountByReferenceID(ctx context.Context, referenceId string) (*entities2.OutStanding, error) {
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

	return &entities2.OutStanding{
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
		go func(i int, loan entities2.Loan) {
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

func (u *BillingUseCase) GetRepaymentInquiryByLoanReferenceId(ctx context.Context, referenceId string) (*entities2.RepaymentInquiry, error) {
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
		return &entities2.RepaymentInquiry{
			LoanId:          loan.Id,
			LoanReferenceId: loan.ReferenceId,
			RepaymentNeeded: nil,
		}, nil
	}

	needRepayments := make([]entities2.RepaymentNeeded, missRepaymentCount)
	for i := 0; i < len(needRepayments); i++ {

		dueDate := loan.CreatedAt.AddDate(0, i, 0)
		isLate := u.Clock.Now().After(dueDate)

		needRepayments[i] = entities2.RepaymentNeeded{
			Amount:  loan.RepaymentAmount,
			DueDate: dueDate,
			IsLate:  isLate,
		}
	}

	return &entities2.RepaymentInquiry{
		LoanId:          loan.Id,
		LoanReferenceId: referenceId,
		RepaymentNeeded: needRepayments,
	}, nil
}

func (u *BillingUseCase) MakePayment(ctx context.Context, repaymentRequest entities2.RepaymentRequest) (int64, error) {
	var (
		errMessage           []string
		loan                 *entities2.Loan
		repaymentTotalAmount int64
		repaymentId          int64

		err error
	)

	if repaymentRequest.LoanReferenceId == "" {
		errMessage = append(errMessage, "loan reference id can not be empty")
	}
	if repaymentRequest.Amount < 1 {
		errMessage = append(errMessage, "amount is invalid")
	}
	if repaymentRequest.RepaymentReferenceId == "" {
		errMessage = append(errMessage, "reference id can not be empty")
	}
	if errMessage != nil || len(errMessage) != 0 {
		return 0, errs.NewWithMessage(http.StatusBadRequest, strings.Join(errMessage, "; "))
	}

	loan, err = u.DBRepo.SelectLoanByReferenceId(ctx, repaymentRequest.LoanReferenceId)
	if err != nil {
		return 0, err
	}

	repaymentTotalAmount, err = u.DBRepo.SelectTotalRepaymentAmountByLoanId(ctx, loan.Id)
	if err != nil {
		if errs.GetHTTPCode(err) != http.StatusNotFound {
			return 0, err
		}
	}

	if repaymentTotalAmount == loan.RepaymentAmount*int64(loan.Tenor) {
		return 0, errs.NewWithMessage(http.StatusBadRequest, "loan has already fully paid")
	}

	if loan.RepaymentAmount != repaymentRequest.Amount {
		return 0, errs.NewWithMessage(http.StatusBadRequest, "payment amount is invalid")
	}

	repaymentId, err = u.DBRepo.CreateRepayment(ctx, entities2.Repayment{
		LoanId:      loan.Id,
		ReferenceId: repaymentRequest.RepaymentReferenceId,
		Amount:      repaymentRequest.Amount,
	})
	if err != nil {
		return 0, err
	}
	return repaymentId, nil
}

func IsUserValid(userId int64) bool {
	if userId < 1 {
		return false
	}

	//TODO: this function has to hit an account service API.
	return false
}
