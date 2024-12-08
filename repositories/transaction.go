package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/sirait-kevin/BillingEngine/domain/entities"
	"github.com/sirait-kevin/BillingEngine/domain/interfaces"
	"github.com/sirait-kevin/BillingEngine/pkg/errs"
)

const (
	insertLoanQuery = `INSERT INTO loans
			(reference_id, user_id, amount, rate_percentage, repayment_amount, status, tenor, repayment_schedule)
			VALUES(?,?,?,?,?,?,?,?);`

	insertRepaymentQuery = `INSERT INTO repayments
			(loan_id, reference_id, amount)
			VALUES(?,?,?);`

	selectLoanByReferenceIdQuery = `SELECT id, reference_id, user_id, amount, rate_percentage, repayment_amount, status, created_at, updated_at, tenor, repayment_schedule
			FROM loans
			WHERE reference_id = ? ORDER BY id DESC;`

	selectActiveLoanByReferenceIdQuery = `SELECT id, reference_id, user_id, amount, rate_percentage, repayment_amount, status, created_at, updated_at, tenor, repayment_schedule
			FROM loans
			WHERE reference_id = ? and status=1;`

	selectLoanByUserIdQuery = `SELECT id, reference_id, user_id, amount, rate_percentage, repayment_amount, status, created_at, updated_at, tenor, repayment_schedule
			FROM loans
			WHERE user_id = ? ORDER BY id DESC;`

	selectRepaymentByReferenceId = `SELECT id, loan_id, reference_id, amount, created_at, updated_at
			FROM repayments
			WHERE reference_id = ?;`

	selectRepaymentByLoanId = `SELECT id, loan_id, reference_id, amount, created_at, updated_at
			FROM repayments
			WHERE loan_id = ? ORDER BY id DESC;`

	selectTotalRepaymentAmountByLoanId = `SELECT IFNULL(SUM(amount), 0)
			FROM repayments
			WHERE loan_id = ?;`

	selectRepaymentCountByLoanId = `SELECT IFNULL(COUNT(id),0)
			FROM repayments
			WHERE loan_id = ?;`

	updateLoanStatusByReferenceId = `UPDATE loans SET status = ? WHERE reference_id = ?;`
)

func (r *DBRepository) CreateLoan(ctx context.Context, tx interfaces.AtomicTransaction, loan entities.Loan) (int64, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("Inserting loan into database: ", loan)
	var (
		err    error
		result sql.Result
	)

	if tx != nil {
		result, err = tx.ExecContext(ctx, insertLoanQuery,
			loan.ReferenceId, loan.UserId, loan.Amount, loan.RatePercentage, loan.RepaymentAmount, loan.Status, loan.Tenor, loan.RepaymentSchedule)
	} else {
		result, err = r.DB.ExecContext(ctx, insertLoanQuery,
			loan.ReferenceId, loan.UserId, loan.Amount, loan.RatePercentage, loan.RepaymentAmount, loan.Status, loan.Tenor, loan.RepaymentSchedule)
	}
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
		loan loansTable
	)

	err = r.DB.GetContext(ctx, &loan, selectLoanByReferenceIdQuery, referenceID)
	if err != nil {
		logger.Error("SelectLoanByReferenceId: ", err)
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
		return nil, err
	}

	return loan.toEntities(), nil

}

func (r *DBRepository) SelectActiveLoanByReferenceId(ctx context.Context, referenceID string) (*entities.Loan, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("select loan by reference id: ", referenceID)
	var (
		err  error
		loan loansTable
	)

	err = r.DB.GetContext(ctx, &loan, selectActiveLoanByReferenceIdQuery, referenceID)
	if err != nil {
		logger.Error("SelectActiveLoanByReferenceId: ", err)
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
		return nil, err
	}

	return loan.toEntities(), nil

}

func (r *DBRepository) SelectLoanByUserId(ctx context.Context, userId int64) (*[]entities.Loan, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("select loan by user id: ", userId)
	var (
		err  error
		loan = []loansTable{}
	)

	err = r.DB.SelectContext(ctx, &loan, selectLoanByUserIdQuery, userId)
	if err != nil {
		logger.Error("SelectLoanByUserId: ", err)
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
		return nil, err
	}

	resp := make([]entities.Loan, len(loan))

	for i, l := range loan {
		resp[i] = *l.toEntities()
	}

	return &resp, nil
}

func (r *DBRepository) CreateRepayment(ctx context.Context, tx interfaces.AtomicTransaction, repayment entities.Repayment) (int64, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("Inserting loan repayment database: ", repayment)
	var (
		err    error
		result sql.Result
	)

	if tx != nil {
		result, err = tx.ExecContext(ctx, insertRepaymentQuery, repayment.LoanId, repayment.ReferenceId, repayment.Amount)
	} else {
		result, err = r.DB.ExecContext(ctx, insertRepaymentQuery,
			repayment.LoanId, repayment.ReferenceId, repayment.Amount)
	}
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
		repayment = repaymentTable{}
	)

	err = r.DB.GetContext(ctx, &repayment, selectRepaymentByReferenceId, referenceID)
	if err != nil {
		logger.Error("SelectRepaymentByReferenceId: ", err)
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
		return nil, err
	}

	return repayment.toEntities(), nil

}

func (r *DBRepository) SelectRepaymentByLoanId(ctx context.Context, loanId int64) (*[]entities.Repayment, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("select repayment by loan id: ", loanId)
	var (
		err        error
		repayments []repaymentTable
	)

	err = r.DB.SelectContext(ctx, &repayments, selectRepaymentByLoanId, loanId)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Error SelectRepaymentByLoanId: ", err)
			err = errs.Wrap(http.StatusNotFound, err)
		}
		return nil, err
	}
	resp := make([]entities.Repayment, len(repayments))
	for i, l := range repayments {
		resp[i] = *l.toEntities()
	}

	return &resp, nil
}

func (r *DBRepository) SelectTotalRepaymentAmountByLoanId(ctx context.Context, loanId int64) (int64, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("select total repayments by loan id: ", loanId)
	var (
		err            error
		totalRepayment int64
	)

	err = r.DB.GetContext(ctx, &totalRepayment, selectTotalRepaymentAmountByLoanId, loanId)
	if err != nil {
		logger.Error("Error fetching SelectTotalRepaymentAmountByLoanId: ", err)
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
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

	err = r.DB.GetContext(ctx, &totalRepayment, selectRepaymentCountByLoanId, loanId)
	if err != nil {
		logger.Error("Error SelectRepaymentCountByLoanId: ", err)
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
		return 0, err
	}

	return totalRepayment, nil
}

func (r *DBRepository) UpdateLoanStatusByReferenceId(ctx context.Context, tx interfaces.AtomicTransaction, referenceId string, status entities.LoanStatus) error {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug(fmt.Sprintf("Update loan status by reference id: %v, status: %v", referenceId, status))

	var (
		err error
		row sql.Result
	)

	if tx != nil {
		row, err = tx.ExecContext(ctx, updateLoanStatusByReferenceId, status, referenceId)
	} else {
		row, err = r.DB.ExecContext(ctx, updateLoanStatusByReferenceId, status, referenceId)
	}

	if err != nil {
		logger.Error("Error UpdateLoanStatusByReferenceId: ", err)
		return err
	}
	if row == nil {
		err = errs.Wrap(http.StatusNotFound, err)
		logger.Error("Error UpdateLoanStatusByReferenceId: ", err)
		return err
	}

	return nil
}
