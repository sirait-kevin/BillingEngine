package repositories

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/sirait-kevin/BillingEngine/entities"
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
)

func (r *UserRepository) CreateLoan(ctx context.Context, loan *entities.Loan) (int64, error) {
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

func (r *UserRepository) CreateRepayment(ctx context.Context, repayment *entities.Repayment) (int64, error) {
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

func (r *UserRepository) SelectRepaymentByReferenceId(ctx context.Context, referenceID string) (*entities.Repayment, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("select repayment by reference id: ", referenceID)
	var (
		err       error
		repayment = &entities.Repayment{}
	)

	err = r.DB.QueryRowContext(ctx, selectLoanByReferenceIdQuery, referenceID).
		Scan(&repayment.Id, &repayment.ReferenceId, &repayment.Amount, &repayment.CreatedAt, &repayment.UpdatedAt)
	if err != nil {
		logger.Error("Error fetching repayment: ", err)
		return nil, err
	}

	return repayment, nil

}

func (r *UserRepository) SelectRepaymentByLoanId(ctx context.Context, referenceID string) (*[]entities.Repayment, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("select repayment by loan id: ", referenceID)
	var (
		err        error
		repayments = &[]entities.Repayment{}
	)

	err = r.DB.SelectContext(ctx, &repayments, selectRepaymentByLoanId)
	if err != nil {
		logger.Error("Error fetching repayment: ", err)
		return nil, err
	}

	return repayments, nil

}
