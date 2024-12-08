package entities

import "time"

type (
	RepaymentInquiry struct {
		LoanId          int64             `json:"loan_id"`
		LoanReferenceId string            `json:"loan_reference_id"`
		LoanStatus      string            `json:"loan_status"`
		RepaymentNeeded []RepaymentNeeded `json:"repayment_needed,omitempty"`
	}

	LoanList struct {
		UserId int64          `json:"user_id"`
		Loans  []LoanResponse `json:"loans,omitempty"`
	}

	LoanResponse struct {
		Id                int64                 `json:"id" `
		ReferenceId       string                `json:"reference_id" `
		UserId            int64                 `json:"user_id,omitempty" `
		Amount            int64                 `json:"amount" `
		RatePercentage    int                   `json:"rate_percentage" `
		Status            string                `json:"status" `
		RepaymentSchedule RepaymentScheduleType `json:"repayment_schedule" `
		Tenor             int                   `json:"tenor" `
		RepaymentAmount   int64                 `json:"repayment_amount" `
		CreatedAt         time.Time             `json:"created_at" `
		UpdatedAt         time.Time             `json:"updated_at,omitempty" `
	}
)
