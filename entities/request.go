package entities

type (
	RepaymentRequest struct {
		LoanReferenceId      string `json:"loan_reference_id"`
		RepaymentReferenceId string `json:"repayment_reference_id"`
		Amount               int64  `json:"amount"`
	}

	LoanRequest struct {
		ReferenceId       string                `json:"reference_id"`
		UserId            int64                 `json:"user_id"`
		Amount            int64                 `json:"amount"`
		RatePercentage    int                   `json:"rate_percentage"`
		RepaymentSchedule RepaymentScheduleType `json:"repayment_schedule"`
		Tenor             int                   `json:"tenor"`
	}
)
