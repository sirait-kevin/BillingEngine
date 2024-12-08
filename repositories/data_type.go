package repositories

import (
	"database/sql"
	"time"

	"github.com/sirait-kevin/BillingEngine/domain/entities"
)

type (
	loansTable struct {
		Id                int64        `db:"id"`
		ReferenceId       string       `db:"reference_id"`
		UserId            int64        `db:"user_id"`
		Amount            int64        `db:"amount"`
		RatePercentage    int          `db:"rate_percentage"`
		Status            int64        `db:"status"`
		RepaymentSchedule string       `db:"repayment_schedule"`
		Tenor             int          `db:"tenor"`
		RepaymentAmount   int64        `db:"repayment_amount"`
		CreatedAt         sql.NullTime `db:"created_at"`
		UpdatedAt         sql.NullTime `db:"updated_at"`
	}

	repaymentTable struct {
		Id          int64        `db:"id"`
		LoanId      int64        `db:"loan_id"`
		ReferenceId string       `db:"reference_id"`
		Amount      int64        `db:"amount"`
		CreatedAt   sql.NullTime `db:"created_at"`
		UpdatedAt   sql.NullTime `db:"updated_at"`
	}
)

func (d *loansTable) toEntities() *entities.Loan {

	var (
		createdAt time.Time
		updatedAt time.Time
	)

	if d.CreatedAt.Valid {
		createdAt = d.CreatedAt.Time
	}
	if d.UpdatedAt.Valid {
		updatedAt = d.UpdatedAt.Time
	}

	return &entities.Loan{
		Id:                d.Id,
		ReferenceId:       d.ReferenceId,
		UserId:            d.UserId,
		Amount:            d.Amount,
		RatePercentage:    d.RatePercentage,
		Status:            entities.LoanStatus(d.Status),
		RepaymentSchedule: entities.RepaymentScheduleType(d.RepaymentSchedule),
		Tenor:             d.Tenor,
		RepaymentAmount:   d.RepaymentAmount,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}

func (d *repaymentTable) toEntities() *entities.Repayment {
	var (
		createdAt time.Time
		updatedAt time.Time
	)

	if d.CreatedAt.Valid {
		createdAt = d.CreatedAt.Time
	}
	if d.UpdatedAt.Valid {
		updatedAt = d.UpdatedAt.Time
	}

	return &entities.Repayment{
		Id:          d.Id,
		LoanId:      d.LoanId,
		ReferenceId: d.ReferenceId,
		Amount:      d.Amount,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
