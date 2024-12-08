package entities

import (
	"strconv"
	"strings"
	"time"

	"github.com/sirait-kevin/BillingEngine/pkg/helper"
)

type (
	Loan struct {
		Id                int64                 `json:"id" `
		ReferenceId       string                `json:"reference_id" `
		UserId            int64                 `json:"user_id" `
		Amount            int64                 `json:"amount" `
		RatePercentage    int                   `json:"rate_percentage" `
		Status            LoanStatus            `json:"status" `
		RepaymentSchedule RepaymentScheduleType `json:"repayment_schedule" `
		Tenor             int                   `json:"tenor" `
		RepaymentAmount   int64                 `json:"repayment_amount" `
		CreatedAt         time.Time             `json:"created_at" `
		UpdatedAt         time.Time             `json:"updated_at" `
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
	param := RepaymentScheduleType(strings.ToLower(string(e)))
	return param == RepaymentMonthly || param == RepaymentWeekly || param == RepaymentYearly
}

func (e LoanStatus) IsActive() bool {
	return e == LoanStatusActive
}

func (e LoanStatus) String() string {
	switch e {
	case LoanStatusActive:
		return "active"
	case LoanStatusRejected:
		return "rejected"
	case LoanStatusCompleted:
		return "completed"
	}
	return "unknown status " + strconv.FormatInt(int64(e), 10)
}

func AddTime(time time.Time, addition int, param RepaymentScheduleType) time.Time {
	switch param {
	case RepaymentMonthly:
		return time.AddDate(0, addition, 0)
	case RepaymentWeekly:
		return time.AddDate(0, 0, addition*7)
	case RepaymentYearly:
		return time.AddDate(addition, 0, 0)
	}
	return time
}

func MissRepayment(createdAt, now time.Time, repaymentCount int, param RepaymentScheduleType) int {
	switch param {
	case RepaymentMonthly:
		return helper.MonthsBetween(createdAt, now) + 1 - repaymentCount
	case RepaymentWeekly:
		return helper.WeeksBetween(createdAt, now) + 1 - repaymentCount
	case RepaymentYearly:
		return helper.YearsBetween(createdAt, now) + 1 - repaymentCount
	}
	return 0
}
