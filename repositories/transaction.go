package repositories

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/sirait-kevin/BillingEngine/entities"
	"github.com/sirait-kevin/BillingEngine/pkg/errs"
)

const (
	insertLoanQuery = `INSERT INTO loans
			(reference_id, user_id, amount, rate_percentage, repayment_amount, status)
			VALUES(?,?,?,?,?,?);`

	insertRepaymentQuery = `INSERT INTO repayments
			(loan_id, reference_id, amount)
			VALUES(?,?,?);`

	selectLoanByReferenceIdQuery = `SELECT id, reference_id, user_id, amount, rate_percentage, repayment_amount, status, created_at, updated_at
			FROM loans
			WHERE reference_id = ?;`

	selectLoanByUserIdQuery = `SELECT id, reference_id, user_id, amount, rate_percentage, repayment_amount, status, created_at, updated_at
			FROM loans
			WHERE user_id = ?;`

	selectRepaymentByReferenceId = `SELECT id, loan_id, reference_id, amount, created_at, updated_at
			FROM repayments
			WHERE reference_id = ?;`

	selectRepaymentByLoanId = `SELECT id, loan_id, reference_id, amount, created_at, updated_at
			FROM repayments
			WHERE loan_id = ?;`

	selectTotalRepaymentAmountByLoanId = `SELECT SUM(amount)
			FROM repayments
			WHERE loan_id = ?;`

	selectRepaymentCountByLoanId = `SELECT COUNT(id)
			FROM repayments
			WHERE loan_id = ?;`
)

func (r *DBRepository) CreateLoan(ctx context.Context, loan entities.Loan) (int64, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("Inserting loan into database: ", loan)
	var (
		err error
	)

	result, err := r.DB.ExecContext(ctx, insertLoanQuery,
		loan.ReferenceId, loan.UserId, loan.Amount, loan.RatePercentage, loan.RepaymentAmount, loan.Status)
	if err != nil {
		logger.Error("Error creating loan: ", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error getting last insert ID: ", err)
		return 0, err
	}
	return id, nil
}

func (r *DBRepository) SelectLoanByReferenceId(ctx context.Context, referenceID string) (*entities.Loan, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("select loan by reference id: ", referenceID)
	var (
		err  error
		loan = &entities.Loan{}
	)

	err = r.DB.GetContext(ctx, &loan, selectLoanByReferenceIdQuery, referenceID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
		logger.Error("Error fetching repayment: ", err)
		return nil, err
	}

	return loan, nil

}

func (r *DBRepository) SelectLoanByUserId(ctx context.Context, userId int64) (*[]entities.Loan, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("select loan by user id: ", userId)
	var (
		err  error
		loan = &[]entities.Loan{}
	)

	err = r.DB.SelectContext(ctx, &loan, selectLoanByUserIdQuery, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
		logger.Error("Error fetching repayment: ", err)
		return nil, err
	}

	return loan, nil
}

func (r *DBRepository) CreateRepayment(ctx context.Context, repayment entities.Repayment) (int64, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("Inserting loan repayment database: ", repayment)
	var (
		err error
	)

	result, err := r.DB.ExecContext(ctx, insertRepaymentQuery,
		repayment.LoanId, repayment.ReferenceId, repayment.Amount)
	if err != nil {
		logger.Error("Error creating loan: ", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error getting last insert ID: ", err)
		return 0, err
	}
	return id, nil
}

func (r *DBRepository) SelectRepaymentByReferenceId(ctx context.Context, referenceID string) (*entities.Repayment, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("select repayment by reference id: ", referenceID)
	var (
		err       error
		repayment = &entities.Repayment{}
	)

	err = r.DB.GetContext(ctx, &repayment, selectLoanByReferenceIdQuery, referenceID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
		logger.Error("Error fetching repayment: ", err)
		return nil, err
	}

	return repayment, nil

}

func (r *DBRepository) SelectRepaymentByLoanId(ctx context.Context, loanId int64) (*[]entities.Repayment, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("select repayment by loan id: ", loanId)
	var (
		err        error
		repayments = &[]entities.Repayment{}
	)

	err = r.DB.SelectContext(ctx, &repayments, selectRepaymentByLoanId, loanId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
		logger.Error("Error fetching repayment: ", err)
		return nil, err
	}

	return repayments, nil
}

func (r *DBRepository) SelectTotalRepaymentAmountByLoanId(ctx context.Context, loanId int64) (int64, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("select total repayments by loan id: ", loanId)
	var (
		err            error
		totalRepayment int64
	)

	err = r.DB.SelectContext(ctx, &totalRepayment, selectTotalRepaymentAmountByLoanId, loanId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
		logger.Error("Error fetching total repayments: ", err)
		return 0, err
	}

	return totalRepayment, nil
}

func (r *DBRepository) SelectRepaymentCountByLoanId(ctx context.Context, loanId int64) (int, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("select repayment count by loan id: ", loanId)
	var (
		err            error
		totalRepayment int
	)

	err = r.DB.SelectContext(ctx, &totalRepayment, selectRepaymentCountByLoanId, loanId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
		logger.Error("Error fetching repayment count: ", err)
		return 0, err
	}

	return totalRepayment, nil
}
