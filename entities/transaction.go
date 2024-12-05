package entities

import "time"

type (
	Loan struct {
		Id                int64     `json:"id"`
		ReferenceId       string    `json:"reference_id"`
		UserId            int64     `json:"user_id"`
		Amount            int64     `json:"amount"`
		RatePercentage    int       `json:"rate_percentage"`
		Status            int       `json:"status"`
		RepaymentSchedule string    `json:"repayment_schedule"`
		Tenor             int       `json:"tenor"`
		RepaymentAmount   int64     `json:"repayment_amount"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
	}

	Repayment struct {
		Id          int64     `json:"id"`
		LoanId      int64     `json:"loan_id"`
		ReferenceId string    `json:"reference_id"`
		Amount      int64     `json:"amount"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	RepaymentHistory struct {
		Loan       Loan        `json:"loan"`
		Repayments []Repayment `json:"repayments"`
	}
)
