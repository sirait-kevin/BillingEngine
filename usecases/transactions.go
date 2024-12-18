package usecases

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/sirait-kevin/BillingEngine/domain/entities"
	"github.com/sirait-kevin/BillingEngine/pkg/errs"
)

func (u *BillingUseCase) CreateLoan(ctx context.Context, loanRequest entities.LoanRequest) (int64, error) {
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

	_, err := u.DBRepo.SelectLoanByReferenceId(ctx, loanRequest.ReferenceId)
	if err == nil {
		return 0, errs.NewWithMessage(http.StatusBadRequest, "reference id already exists")
	}
	if errs.GetHTTPCode(err) != http.StatusNotFound {
		return 0, err
	}

	isDelinquent, err := u.GetUserStatusIsDelinquent(ctx, loanRequest.UserId)
	if err != nil {
		return 0, err
	}
	if isDelinquent {
		return 0, errs.NewWithMessage(http.StatusForbidden, "User is delinquent")
	}

	repaymentAmount := (loanRequest.Amount + (loanRequest.Amount * int64(loanRequest.RatePercentage) / 100)) / int64(loanRequest.Tenor)

	loanId, err := u.DBRepo.CreateLoan(ctx, nil, entities.Loan{
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

func (u *BillingUseCase) GetPaymentHistoryByReferenceID(ctx context.Context, referenceId string) (*entities.LoanHistory, error) {

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
	if !IsUserValid(userId) {
		return false, errs.NewWithMessage(http.StatusBadRequest, "user id is invalid")
	}

	loans, err := u.DBRepo.SelectLoanByUserId(ctx, userId)
	if err != nil {
		if errs.GetHTTPCode(err) != http.StatusNotFound {
			return false, err
		} else {
			return false, nil
		}
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
				if errs.GetHTTPCode(errWg[i]) == http.StatusNotFound {
					errWg[i] = nil
				}
			}
		}(i, loan)
	}
	wg.Wait()

	for _, errW := range errWg {
		if errW != nil {
			return false, errW
		}
	}

	for _, loan := range *loans {
		if loan.Tenor > repaymentCounts[loan.Id] {
			if entities.MissRepayment(loan.CreatedAt, u.Clock.Now(), repaymentCounts[loan.Id], loan.RepaymentSchedule) > 1 {
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

	missRepaymentCount := entities.MissRepayment(loan.CreatedAt, u.Clock.Now(), repaymentCount, loan.RepaymentSchedule)

	needRepayments := []entities.RepaymentNeeded{}
	for i := 0; i < missRepaymentCount; i++ {
		addTime := i + repaymentCount + 1
		dueDate := entities.AddTime(loan.CreatedAt, addTime, loan.RepaymentSchedule)
		isLate := u.Clock.Now().After(dueDate)

		needRepayments = append(needRepayments, entities.RepaymentNeeded{
			Amount:  loan.RepaymentAmount,
			DueDate: dueDate,
			IsLate:  isLate,
		})
	}

	if len(needRepayments) == 0 && loan.Status.IsActive() {
		addTime := repaymentCount + 1
		dueDate := entities.AddTime(loan.CreatedAt, addTime, loan.RepaymentSchedule)
		isLate := u.Clock.Now().After(dueDate)
		needRepayments = append(needRepayments, entities.RepaymentNeeded{
			Amount:  loan.RepaymentAmount,
			DueDate: dueDate,
			IsLate:  isLate,
		})
	}

	return &entities.RepaymentInquiry{
		LoanId:          loan.Id,
		LoanReferenceId: loan.ReferenceId,
		LoanStatus:      loan.Status.String(),
		RepaymentNeeded: needRepayments,
	}, nil
}

func (u *BillingUseCase) MakePayment(ctx context.Context, repaymentRequest entities.RepaymentRequest) (int64, error) {
	var (
		errMessage           []string
		loan                 *entities.Loan
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

	_, err = u.DBRepo.SelectRepaymentByReferenceId(ctx, repaymentRequest.RepaymentReferenceId)
	if err == nil {
		return 0, errs.NewWithMessage(http.StatusBadRequest, "reference id already exists")
	}
	if errs.GetHTTPCode(err) != http.StatusNotFound {
		return 0, err
	}

	loan, err = u.DBRepo.SelectLoanByReferenceId(ctx, repaymentRequest.LoanReferenceId)
	if err != nil {
		return 0, err
	}
	if !loan.Status.IsActive() {
		return 0, errs.NewWithMessage(http.StatusBadRequest, "loan status has been "+loan.Status.String())
	}

	repaymentTotalAmount, err = u.DBRepo.SelectTotalRepaymentAmountByLoanId(ctx, loan.Id)
	if err != nil {
		if errs.GetHTTPCode(err) != http.StatusNotFound {
			return 0, err
		}
	}

	if loan.RepaymentAmount != repaymentRequest.Amount {
		return 0, errs.NewWithMessage(http.StatusBadRequest, "payment amount is invalid, expected: "+strconv.FormatInt(loan.RepaymentAmount, 10))
	}

	dbTx, err := u.DBRepo.BeginTx(ctx)
	if err != nil {
		return 0, err
	}
	defer dbTx.Rollback()

	repaymentId, err = u.DBRepo.CreateRepayment(ctx, dbTx, entities.Repayment{
		LoanId:      loan.Id,
		ReferenceId: repaymentRequest.RepaymentReferenceId,
		Amount:      repaymentRequest.Amount,
	})
	if err != nil {
		return 0, err
	}

	if repaymentTotalAmount+repaymentRequest.Amount >= loan.RepaymentAmount*int64(loan.Tenor) {
		err = u.DBRepo.UpdateLoanStatusByReferenceId(ctx, dbTx, loan.ReferenceId, entities.LoanStatusCompleted)
		if err != nil {
			return 0, err
		}
	}

	err = dbTx.Commit()
	if err != nil {
		return 0, err
	}

	return repaymentId, nil
}

func (u *BillingUseCase) GetLoanListByUserId(ctx context.Context, userId int64) (*[]entities.Loan, error) {
	if userId < 1 {
		return nil, errs.NewWithMessage(http.StatusBadRequest, "user id can not be empty")
	}

	loans, err := u.DBRepo.SelectLoanByUserId(ctx, userId)
	if err != nil {
		if errs.GetHTTPCode(err) != http.StatusNotFound {
			return nil, err
		}
	}

	return loans, err
}

func IsUserValid(userId int64) bool {
	if userId < 1 {
		return false
	}

	//TODO: this function has to hit an account service API.
	return true
}
