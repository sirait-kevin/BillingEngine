package restful

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirait-kevin/BillingEngine/domain/entities"
	"github.com/sirait-kevin/BillingEngine/pkg/errs"
	"github.com/sirait-kevin/BillingEngine/pkg/helper"
)

func (h *BillingHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var loanRequest entities.LoanRequest
	err := json.NewDecoder(r.Body).Decode(&loanRequest)
	if err != nil {
		helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	loanId, err := h.BillingUC.CreateLoan(ctx, loanRequest)
	if err != nil {
		helper.JSON(w, ctx, nil, err)
		return
	}

	helper.JSON(w, ctx, map[string]int64{
		"loan_id": loanId,
	}, nil)
}

func (h *BillingHandler) GetPaymentHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reference_id := r.FormValue("reference_id")

	paymentHistory, err := h.BillingUC.GetPaymentHistoryByReferenceID(ctx, reference_id)
	if err != nil {
		helper.JSON(w, ctx, nil, err)
		return
	}

	helper.JSON(w, ctx, paymentHistory, nil)
}

func (h *BillingHandler) MakePayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var paymentRequest entities.RepaymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	paymentID, err := h.BillingUC.MakePayment(ctx, paymentRequest)
	if err != nil {
		helper.JSON(w, ctx, nil, err)
		return
	}

	helper.JSON(w, ctx, map[string]int64{
		"payment_id": paymentID,
	}, nil)

}

func (h *BillingHandler) GetOutStandingAmount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	referenceId := r.FormValue("reference_id")

	inquiry, err := h.BillingUC.GetOutStandingAmountByReferenceID(ctx, referenceId)
	if err != nil {
		helper.JSON(w, ctx, nil, err)
		return
	}

	helper.JSON(w, ctx, inquiry, nil)
}

func (h *BillingHandler) GetUserStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	if err != nil {
		helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusBadRequest, "Invalid user ID"))
		return
	}
	isDelinquent, err := h.BillingUC.GetUserStatusIsDelinquent(ctx, userId)
	if err != nil {
		helper.JSON(w, ctx, nil, err)
		return
	}
	helper.JSON(w, ctx, map[string]bool{
		"is_delinquent": isDelinquent,
	}, nil)
}

func (h *BillingHandler) GetPaymentInquiry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	referenceId := r.FormValue("loan_reference_id")

	repaymentInquiry, err := h.BillingUC.GetRepaymentInquiryByLoanReferenceId(ctx, referenceId)
	if err != nil {
		helper.JSON(w, ctx, nil, err)
		return
	}

	helper.JSON(w, ctx, repaymentInquiry, nil)
}

func (h *BillingHandler) GetLoanHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, err := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	if err != nil {
		helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	loans, err := h.BillingUC.GetLoanListByUserId(ctx, userId)
	if err != nil {
		helper.JSON(w, ctx, nil, err)
		return
	}

	loanResponses := make([]entities.LoanResponse, len(*loans))
	for i, l := range *loans {
		loanResponses[i] = entities.LoanResponse{
			Id:                l.Id,
			ReferenceId:       l.ReferenceId,
			Amount:            l.Amount,
			RatePercentage:    l.RatePercentage,
			Status:            l.Status.String(),
			RepaymentSchedule: l.RepaymentSchedule,
			Tenor:             l.Tenor,
			RepaymentAmount:   l.RepaymentAmount,
			CreatedAt:         l.CreatedAt,
			UpdatedAt:         l.UpdatedAt,
		}
	}

	helper.JSON(w, ctx, &entities.LoanList{
		UserId: userId,
		Loans:  loanResponses,
	}, nil)

}
