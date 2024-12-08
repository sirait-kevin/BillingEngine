package restful

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sirait-kevin/BillingEngine/domain/entities"
	mock_handler "github.com/sirait-kevin/BillingEngine/mocks/handler"
)

func getSampleCreateLoanRequest() entities.LoanRequest {
	return entities.LoanRequest{
		ReferenceId:       "",
		UserId:            1,
		Amount:            1,
		RatePercentage:    1,
		RepaymentSchedule: "",
		Tenor:             1,
	}

}

func getSampleMakePaymentRequest() entities.RepaymentRequest {
	return entities.RepaymentRequest{}
}

func TestBillingHandler_CreateLoan(t *testing.T) {
	type fields struct {
		BillingUC *mock_handler.MockBillingUsecase
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   func(ctrl *gomock.Controller) fields
		args     args
		mock     func(f fields, args args)
		wantCode int
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					reqBody := getSampleCreateLoanRequest()
					jsonB, _ := json.Marshal(reqBody)

					r := httptest.NewRequest("POST", "localhost:8080/create/loan", bytes.NewBuffer(jsonB))
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().CreateLoan(gomock.Any(), getSampleCreateLoanRequest()).Return(int64(1), nil)
			},
			wantCode: 200,
		},
		{
			name: "error request decoding",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					reqBody := "error"
					jsonB := []byte(reqBody)
					r := httptest.NewRequest("POST", "localhost:8080/create/loan", bytes.NewBuffer(jsonB))
					return r
				}(),
			},
			mock: func(f fields, args args) {
			},
			wantCode: 400,
		},
		{
			name: "error usecase",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					reqBody := getSampleCreateLoanRequest()
					jsonB, _ := json.Marshal(reqBody)

					r := httptest.NewRequest("POST", "localhost:8080/create/loan", bytes.NewBuffer(jsonB))
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().CreateLoan(gomock.Any(), getSampleCreateLoanRequest()).Return(int64(1), errors.New("some error"))
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			h := &BillingHandler{
				BillingUC: f.BillingUC,
			}
			tt.mock(f, tt.args)

			h.CreateLoan(tt.args.w, tt.args.r)
			assert.EqualValues(t, tt.wantCode, tt.args.w.Code)
		})
	}
}

func TestBillingHandler_GetPaymentHistory(t *testing.T) {
	type fields struct {
		BillingUC *mock_handler.MockBillingUsecase
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   func(ctrl *gomock.Controller) fields
		args     args
		mock     func(f fields, args args)
		wantCode int
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "localhost:8080/payment/history", nil)
					r.Form = url.Values{
						"reference_id": {"reference"},
					}
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().GetPaymentHistoryByReferenceID(gomock.Any(), args.r.Form.Get("reference_id")).Return(&entities.LoanHistory{}, nil)
			},
			wantCode: 200,
		},
		{
			name: "error usecase",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "localhost:8080/payment/history", nil)
					r.Form = url.Values{
						"reference_id": {"reference"},
					}
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().GetPaymentHistoryByReferenceID(gomock.Any(), args.r.Form.Get("reference_id")).Return(&entities.LoanHistory{}, errors.New("some error"))
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			h := &BillingHandler{
				BillingUC: f.BillingUC,
			}
			tt.mock(f, tt.args)

			h.GetPaymentHistory(tt.args.w, tt.args.r)
			assert.EqualValues(t, tt.wantCode, tt.args.w.Code)
		})
	}
}

func TestBillingHandler_MakePayment(t *testing.T) {
	type fields struct {
		BillingUC *mock_handler.MockBillingUsecase
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   func(ctrl *gomock.Controller) fields
		args     args
		mock     func(f fields, args args)
		wantCode int
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					reqBody := getSampleMakePaymentRequest()
					jsonB, _ := json.Marshal(reqBody)

					r := httptest.NewRequest("POST", "localhost:8080/make/payment", bytes.NewBuffer(jsonB))
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().MakePayment(gomock.Any(), getSampleMakePaymentRequest()).Return(int64(1), nil)
			},
			wantCode: 200,
		},
		{
			name: "error request decoding",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					reqBody := "error"
					jsonB, _ := json.Marshal(reqBody)

					r := httptest.NewRequest("POST", "localhost:8080/make/payment", bytes.NewBuffer(jsonB))
					return r
				}(),
			},
			mock: func(f fields, args args) {
			},
			wantCode: 400,
		},
		{
			name: "error usecase",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					reqBody := getSampleMakePaymentRequest()
					jsonB, _ := json.Marshal(reqBody)

					r := httptest.NewRequest("POST", "localhost:8080/make/payment", bytes.NewBuffer(jsonB))
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().MakePayment(gomock.Any(), getSampleMakePaymentRequest()).Return(int64(1), errors.New("some error"))
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			h := &BillingHandler{
				BillingUC: f.BillingUC,
			}
			tt.mock(f, tt.args)

			h.MakePayment(tt.args.w, tt.args.r)
			assert.EqualValues(t, tt.wantCode, tt.args.w.Code)
		})
	}
}

func TestBillingHandler_GetOutStandingAmount(t *testing.T) {
	type fields struct {
		BillingUC *mock_handler.MockBillingUsecase
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   func(ctrl *gomock.Controller) fields
		args     args
		mock     func(f fields, args args)
		wantCode int
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "localhost:8080/outstanding/amount", nil)
					r.Form = url.Values{
						"reference_id": {"reference"},
					}
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().GetOutStandingAmountByReferenceID(gomock.Any(), args.r.Form.Get("reference_id")).Return(&entities.OutStanding{}, nil)
			},
			wantCode: 200,
		},
		{
			name: "error usecase",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "localhost:8080/outstanding/amount", nil)
					r.Form = url.Values{
						"reference_id": {"reference"},
					}
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().GetOutStandingAmountByReferenceID(gomock.Any(), args.r.Form.Get("reference_id")).Return(&entities.OutStanding{}, errors.New("some error"))
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			h := &BillingHandler{
				BillingUC: f.BillingUC,
			}
			tt.mock(f, tt.args)

			h.GetOutStandingAmount(tt.args.w, tt.args.r)
			assert.EqualValues(t, tt.wantCode, tt.args.w.Code)
		})
	}
}

func TestBillingHandler_GetUserStatus(t *testing.T) {
	type fields struct {
		BillingUC *mock_handler.MockBillingUsecase
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   func(ctrl *gomock.Controller) fields
		args     args
		mock     func(f fields, args args)
		wantCode int
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "localhost:8080/user/status", nil)
					r.Form = url.Values{
						"user_id": {"1"},
					}
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().GetUserStatusIsDelinquent(gomock.Any(), int64(1)).Return(true, nil)
			},
			wantCode: 200,
		},
		{
			name: "error parameter",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "localhost:8080/user/status", nil)
					r.Form = url.Values{
						"user_id": {"user"},
					}
					return r
				}(),
			},
			mock: func(f fields, args args) {
			},
			wantCode: 400,
		},
		{
			name: "error usecase",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "localhost:8080/user/status", nil)
					r.Form = url.Values{
						"user_id": {"1"},
					}
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().GetUserStatusIsDelinquent(gomock.Any(), int64(1)).Return(false, errors.New("some error"))
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			h := &BillingHandler{
				BillingUC: f.BillingUC,
			}
			tt.mock(f, tt.args)

			h.GetUserStatus(tt.args.w, tt.args.r)
			assert.EqualValues(t, tt.wantCode, tt.args.w.Code)
		})
	}
}

func TestBillingHandler_GetPaymentInquiry(t *testing.T) {
	type fields struct {
		BillingUC *mock_handler.MockBillingUsecase
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   func(ctrl *gomock.Controller) fields
		args     args
		mock     func(f fields, args args)
		wantCode int
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "localhost:8080/payment/inquiry", nil)
					r.Form = url.Values{
						"loan_reference_id": {"reference"},
					}
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().GetRepaymentInquiryByLoanReferenceId(gomock.Any(), args.r.Form.Get("loan_reference_id")).Return(&entities.RepaymentInquiry{}, nil)
			},
			wantCode: 200,
		},
		{
			name: "error usecase",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "localhost:8080/payment/inquiry", nil)
					r.Form = url.Values{
						"loan_reference_id": {"reference"},
					}
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().GetRepaymentInquiryByLoanReferenceId(gomock.Any(), args.r.Form.Get("loan_reference_id")).Return(&entities.RepaymentInquiry{}, errors.New("some error"))
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			h := &BillingHandler{
				BillingUC: f.BillingUC,
			}
			tt.mock(f, tt.args)

			h.GetPaymentInquiry(tt.args.w, tt.args.r)
			assert.EqualValues(t, tt.wantCode, tt.args.w.Code)
		})
	}
}

func TestBillingHandler_GetLoanHistory(t *testing.T) {
	type fields struct {
		BillingUC *mock_handler.MockBillingUsecase
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   func(ctrl *gomock.Controller) fields
		args     args
		mock     func(f fields, args args)
		wantCode int
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "localhost:8080/loan/history", nil)
					r.Form = url.Values{
						"user_id": {"1"},
					}
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().GetLoanListByUserId(gomock.Any(), int64(1)).Return(&[]entities.Loan{
					{
						Id:                0,
						ReferenceId:       "",
						UserId:            0,
						Amount:            0,
						RatePercentage:    0,
						Status:            0,
						RepaymentSchedule: "",
						Tenor:             0,
						RepaymentAmount:   0,
						CreatedAt:         time.Time{},
						UpdatedAt:         time.Time{},
					},
					{
						Id:                0,
						ReferenceId:       "",
						UserId:            0,
						Amount:            0,
						RatePercentage:    0,
						Status:            0,
						RepaymentSchedule: "",
						Tenor:             0,
						RepaymentAmount:   0,
						CreatedAt:         time.Time{},
						UpdatedAt:         time.Time{},
					},
				}, nil)
			},
			wantCode: 200,
		},
		{
			name: "error parameter",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "localhost:8080/loan/history", nil)
					r.Form = url.Values{
						"user_id": {"user"},
					}
					return r
				}(),
			},
			mock: func(f fields, args args) {
			},
			wantCode: 400,
		},
		{
			name: "error usecase",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					BillingUC: mock_handler.NewMockBillingUsecase(ctrl),
				}
			},
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					r := httptest.NewRequest("GET", "localhost:8080/loan/history", nil)
					r.Form = url.Values{
						"user_id": {"1"},
					}
					return r
				}(),
			},
			mock: func(f fields, args args) {
				f.BillingUC.EXPECT().GetLoanListByUserId(gomock.Any(), int64(1)).Return(&[]entities.Loan{}, errors.New("some error"))
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			h := &BillingHandler{
				BillingUC: f.BillingUC,
			}
			tt.mock(f, tt.args)

			h.GetLoanHistory(tt.args.w, tt.args.r)
			assert.EqualValues(t, tt.wantCode, tt.args.w.Code)
		})
	}
}
