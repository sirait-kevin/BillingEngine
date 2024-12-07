package entities

import "time"

type (
	Loan struct {
		Id                int64                 `json:"id"`
		ReferenceId       string                `json:"reference_id"`
		UserId            int64                 `json:"user_id"`
		Amount            int64                 `json:"amount"`
		RatePercentage    int                   `json:"rate_percentage"`
		Status            LoanStatus            `json:"status"`
		RepaymentSchedule RepaymentScheduleType `json:"repayment_schedule"`
		Tenor             int                   `json:"tenor"`
		RepaymentAmount   int64                 `json:"repayment_amount"`
		CreatedAt         time.Time             `json:"created_at"`
		UpdatedAt         time.Time             `json:"updated_at"`
	}

	Repayment struct {
		Id          int64     `json:"id"`
		LoanId      int64     `json:"loan_id"`
		ReferenceId string    `json:"reference_id"`
		Amount      int64     `json:"amount"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	LoanHistory struct {
		Loan       Loan        `json:"loan"`
		Repayments []Repayment `json:"repayments"`
	}

	OutStanding struct {
		LoanId            int64  `json:"loan_id"`
		LoanReferenceId   string `json:"loan_reference_id"`
		OutstandingAmount int64  `json:"outstanding_amount"`
	}

	RepaymentInquiry struct {
		LoanId          int64             `json:"loan_id"`
		LoanReferenceId string            `json:"loan_reference_id"`
		RepaymentNeeded []RepaymentNeeded `json:"repayment_needed"`
	}

	RepaymentNeeded struct {
		Amount  int64     `json:"amount"`
		DueDate time.Time `json:"due_date"`
		IsLate  bool      `json:"is_late"`
	}

	LoanStatus            int
	RepaymentScheduleType string
)

const (
	LoanStatusActive    LoanStatus = 1
	LoanStatusRejected  LoanStatus = 2
	LoanStatusCompleted LoanStatus = 3

	RepaymentMonthly RepaymentScheduleType = "monthly"
	RepaymentWeekly  RepaymentScheduleType = "weekly"
	RepaymentYearly  RepaymentScheduleType = "yearly"
)

func (e RepaymentScheduleType) IsValid() bool {
	return e == RepaymentMonthly || e == RepaymentWeekly || e == RepaymentYearly
}
