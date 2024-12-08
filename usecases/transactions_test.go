package usecases

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/sirait-kevin/BillingEngine/domain/entities"
	mock_domain "github.com/sirait-kevin/BillingEngine/mocks/domain"
	mock_usecase "github.com/sirait-kevin/BillingEngine/mocks/usecases"
	"github.com/sirait-kevin/BillingEngine/pkg/errs"
)

func TestBillingUseCase_CreateLoan(t *testing.T) {
	type input struct {
		ctx   context.Context
		param entities.LoanRequest
	}
	type fields struct {
		DBRepo *mock_usecase.MockDBRepository
		Clock  *mock_domain.MockClock
	}
	tests := []struct {
		name    string
		fields  func(ctrl *gomock.Controller) fields
		input   input
		mock    func(f fields, input input)
		want    int64
		wantErr bool
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.LoanRequest{
					ReferenceId:       "1",
					UserId:            1,
					Amount:            1,
					RatePercentage:    1,
					RepaymentSchedule: "weekly",
					Tenor:             1,
				},
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.ReferenceId).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
				f.DBRepo.EXPECT().SelectLoanByUserId(gomock.Any(), args.param.UserId).Return(&[]entities.Loan{}, nil)

				repaymentAmount := (args.param.Amount + (args.param.Amount * int64(args.param.RatePercentage) / 100)) / int64(args.param.Tenor)
				f.DBRepo.EXPECT().CreateLoan(gomock.Any(), nil, entities.Loan{
					ReferenceId:       args.param.ReferenceId,
					UserId:            args.param.UserId,
					Amount:            args.param.Amount,
					RatePercentage:    args.param.RatePercentage,
					Status:            entities.LoanStatusActive,
					RepaymentSchedule: args.param.RepaymentSchedule,
					Tenor:             args.param.Tenor,
					RepaymentAmount:   repaymentAmount,
				}).Return(int64(1), nil)

			},
			want:    1,
			wantErr: false,
		},
		{
			name: "error parameter",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.LoanRequest{
					RatePercentage: -1,
				},
			},
			mock: func(f fields, args input) {
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error already exist",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.LoanRequest{
					ReferenceId:       "1",
					UserId:            1,
					Amount:            1,
					RatePercentage:    1,
					RepaymentSchedule: "weekly",
					Tenor:             1,
				},
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.ReferenceId).Return(nil, nil)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error select loan",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.LoanRequest{
					ReferenceId:       "1",
					UserId:            1,
					Amount:            1,
					RatePercentage:    1,
					RepaymentSchedule: "weekly",
					Tenor:             1,
				},
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.ReferenceId).Return(nil, errs.NewWithMessage(http.StatusForbidden, ""))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error get is delinquent",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.LoanRequest{
					ReferenceId:       "1",
					UserId:            1,
					Amount:            1,
					RatePercentage:    1,
					RepaymentSchedule: "weekly",
					Tenor:             1,
				},
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.ReferenceId).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
				f.DBRepo.EXPECT().SelectLoanByUserId(gomock.Any(), args.param.UserId).Return(&[]entities.Loan{}, errs.NewWithMessage(http.StatusForbidden, ""))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error is deliquent",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.LoanRequest{
					ReferenceId:       "1",
					UserId:            1,
					Amount:            1,
					RatePercentage:    1,
					RepaymentSchedule: "weekly",
					Tenor:             1,
				},
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.ReferenceId).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
				f.DBRepo.EXPECT().SelectLoanByUserId(gomock.Any(), args.param.UserId).Return(&[]entities.Loan{
					{
						Id:                1,
						ReferenceId:       "",
						UserId:            0,
						Amount:            1000,
						RatePercentage:    0,
						Status:            0,
						RepaymentSchedule: entities.RepaymentMonthly,
						Tenor:             2,
						RepaymentAmount:   0,
						CreatedAt:         time.Date(2000, 10, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:         time.Time{},
					},
				}, nil)
				f.DBRepo.EXPECT().SelectRepaymentCountByLoanId(gomock.Any(), int64(1)).Return(int(0), nil)
				f.Clock.EXPECT().Now().Return(time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error create loan",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.LoanRequest{
					ReferenceId:       "1",
					UserId:            1,
					Amount:            1,
					RatePercentage:    1,
					RepaymentSchedule: "weekly",
					Tenor:             1,
				},
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.ReferenceId).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
				f.DBRepo.EXPECT().SelectLoanByUserId(gomock.Any(), args.param.UserId).Return(&[]entities.Loan{}, nil)

				repaymentAmount := (args.param.Amount + (args.param.Amount * int64(args.param.RatePercentage) / 100)) / int64(args.param.Tenor)
				f.DBRepo.EXPECT().CreateLoan(gomock.Any(), nil, entities.Loan{
					ReferenceId:       args.param.ReferenceId,
					UserId:            args.param.UserId,
					Amount:            args.param.Amount,
					RatePercentage:    args.param.RatePercentage,
					Status:            entities.LoanStatusActive,
					RepaymentSchedule: args.param.RepaymentSchedule,
					Tenor:             args.param.Tenor,
					RepaymentAmount:   repaymentAmount,
				}).Return(int64(1), errs.NewWithMessage(http.StatusInternalServerError, ""))

			},
			want:    1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			u := BillingUseCase{
				DBRepo: f.DBRepo,
				Clock:  f.Clock,
			}
			tt.mock(f, tt.input)

			got, err := u.CreateLoan(tt.input.ctx, tt.input.param)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestBillingUseCase_GetPaymentHistoryByReferenceID(t *testing.T) {
	type input struct {
		ctx   context.Context
		param string
	}
	type fields struct {
		DBRepo *mock_usecase.MockDBRepository
		Clock  *mock_domain.MockClock
	}
	tests := []struct {
		name    string
		fields  func(ctrl *gomock.Controller) fields
		input   input
		mock    func(f fields, input input)
		want    *entities.LoanHistory
		wantErr bool
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "reference",
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param).Return(&entities.Loan{
					Id:                1,
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
				}, nil)
				f.DBRepo.EXPECT().SelectRepaymentByLoanId(gomock.Any(), int64(1)).Return(&[]entities.Repayment{
					{
						Id:          0,
						LoanId:      0,
						ReferenceId: "",
						Amount:      0,
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
					},
				}, nil)

			},
			want: &entities.LoanHistory{
				Loan: entities.Loan{
					Id: 1,
				},
				Repayments: []entities.Repayment{
					{},
				},
			},
			wantErr: false,
		},
		{
			name: "error parameter",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "",
			},
			mock: func(f fields, args input) {
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "error select loan",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "reference",
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param).Return(&entities.Loan{
					Id:                1,
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
				}, errs.NewWithMessage(http.StatusInternalServerError, ""))

			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error select repayment",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "reference",
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param).Return(&entities.Loan{
					Id:                1,
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
				}, nil)
				f.DBRepo.EXPECT().SelectRepaymentByLoanId(gomock.Any(), int64(1)).Return(&[]entities.Repayment{
					{
						Id:          0,
						LoanId:      0,
						ReferenceId: "",
						Amount:      0,
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
					},
				}, errors.New("SelectRepaymentByLoanIdError"))

			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			u := BillingUseCase{
				DBRepo: f.DBRepo,
				Clock:  f.Clock,
			}
			tt.mock(f, tt.input)

			got, err := u.GetPaymentHistoryByReferenceID(tt.input.ctx, tt.input.param)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestBillingUseCase_GetOutStandingAmountByReferenceID(t *testing.T) {
	type input struct {
		ctx   context.Context
		param string
	}
	type fields struct {
		DBRepo *mock_usecase.MockDBRepository
		Clock  *mock_domain.MockClock
	}
	tests := []struct {
		name    string
		fields  func(ctrl *gomock.Controller) fields
		input   input
		mock    func(f fields, input input)
		want    *entities.OutStanding
		wantErr bool
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "reference",
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            1000,
					RatePercentage:    0,
					Status:            0,
					RepaymentSchedule: "",
					Tenor:             1,
					RepaymentAmount:   1000,
					CreatedAt:         time.Time{},
					UpdatedAt:         time.Time{},
				}, nil)
				f.DBRepo.EXPECT().SelectTotalRepaymentAmountByLoanId(gomock.Any(), int64(1)).Return(int64(0), errs.NewWithMessage(http.StatusNotFound, ""))

			},
			want: &entities.OutStanding{
				LoanId:            1,
				LoanReferenceId:   "",
				OutstandingAmount: 1000,
			},
			wantErr: false,
		},
		{
			name: "error parameter",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "",
			},
			mock: func(f fields, args input) {
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error select loan",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "reference",
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            1000,
					RatePercentage:    0,
					Status:            0,
					RepaymentSchedule: "",
					Tenor:             1,
					RepaymentAmount:   1000,
					CreatedAt:         time.Time{},
					UpdatedAt:         time.Time{},
				}, errs.NewWithMessage(http.StatusInternalServerError, ""))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error select total amount",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "reference",
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            1000,
					RatePercentage:    0,
					Status:            0,
					RepaymentSchedule: "",
					Tenor:             1,
					RepaymentAmount:   1000,
					CreatedAt:         time.Time{},
					UpdatedAt:         time.Time{},
				}, nil)
				f.DBRepo.EXPECT().SelectTotalRepaymentAmountByLoanId(gomock.Any(), int64(1)).Return(int64(0), errs.NewWithMessage(http.StatusInternalServerError, ""))

			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			u := BillingUseCase{
				DBRepo: f.DBRepo,
				Clock:  f.Clock,
			}
			tt.mock(f, tt.input)

			got, err := u.GetOutStandingAmountByReferenceID(tt.input.ctx, tt.input.param)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestBillingUseCase_GetUserStatusIsDelinquent(t *testing.T) {
	type input struct {
		ctx   context.Context
		param int64
	}
	type fields struct {
		DBRepo *mock_usecase.MockDBRepository
		Clock  *mock_domain.MockClock
	}
	tests := []struct {
		name    string
		fields  func(ctrl *gomock.Controller) fields
		input   input
		mock    func(f fields, input input)
		want    bool
		wantErr bool
	}{
		{
			name: "success not delinquent",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: 1,
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByUserId(gomock.Any(), args.param).Return(&[]entities.Loan{
					{
						Id:                1,
						ReferenceId:       "",
						UserId:            0,
						Amount:            1000,
						RatePercentage:    0,
						Status:            0,
						RepaymentSchedule: entities.RepaymentMonthly,
						Tenor:             2,
						RepaymentAmount:   0,
						CreatedAt:         time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:         time.Time{},
					},
				}, nil)
				f.DBRepo.EXPECT().SelectRepaymentCountByLoanId(gomock.Any(), int64(1)).Return(int(0), nil)
				f.Clock.EXPECT().Now().Return(time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC))
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "error parameter",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: 0,
			},
			mock: func(f fields, args input) {
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "error select loan not found",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: 1,
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByUserId(gomock.Any(), args.param).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "error select loan",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: 1,
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByUserId(gomock.Any(), args.param).Return(&[]entities.Loan{
					{
						Id:                1,
						ReferenceId:       "",
						UserId:            0,
						Amount:            1000,
						RatePercentage:    0,
						Status:            0,
						RepaymentSchedule: entities.RepaymentMonthly,
						Tenor:             2,
						RepaymentAmount:   0,
						CreatedAt:         time.Date(2000, 10, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:         time.Time{},
					},
				}, errs.NewWithMessage(http.StatusInternalServerError, ""))
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "error select count",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: 1,
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByUserId(gomock.Any(), args.param).Return(&[]entities.Loan{
					{
						Id:                1,
						ReferenceId:       "",
						UserId:            0,
						Amount:            1000,
						RatePercentage:    0,
						Status:            0,
						RepaymentSchedule: entities.RepaymentMonthly,
						Tenor:             2,
						RepaymentAmount:   0,
						CreatedAt:         time.Date(2000, 10, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:         time.Time{},
					},
				}, nil)
				f.DBRepo.EXPECT().SelectRepaymentCountByLoanId(gomock.Any(), int64(1)).Return(int(0), errs.NewWithMessage(http.StatusInternalServerError, ""))
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "success delinquent",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: 1,
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByUserId(gomock.Any(), args.param).Return(&[]entities.Loan{
					{
						Id:                1,
						ReferenceId:       "",
						UserId:            0,
						Amount:            1000,
						RatePercentage:    0,
						Status:            0,
						RepaymentSchedule: entities.RepaymentMonthly,
						Tenor:             2,
						RepaymentAmount:   0,
						CreatedAt:         time.Date(2000, 10, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:         time.Time{},
					},
				}, nil)
				f.DBRepo.EXPECT().SelectRepaymentCountByLoanId(gomock.Any(), int64(1)).Return(int(0), errs.NewWithMessage(http.StatusNotFound, ""))
				f.Clock.EXPECT().Now().Return(time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC))
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			u := BillingUseCase{
				DBRepo: f.DBRepo,
				Clock:  f.Clock,
			}
			tt.mock(f, tt.input)

			got, err := u.GetUserStatusIsDelinquent(tt.input.ctx, tt.input.param)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestBillingUseCase_GetRepaymentInquiryByLoanReferenceId(t *testing.T) {
	type input struct {
		ctx   context.Context
		param string
	}
	type fields struct {
		DBRepo *mock_usecase.MockDBRepository
		Clock  *mock_domain.MockClock
	}
	tests := []struct {
		name    string
		fields  func(ctrl *gomock.Controller) fields
		input   input
		mock    func(f fields, input input)
		want    *entities.RepaymentInquiry
		wantErr bool
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "reference",
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            1000,
					RatePercentage:    0,
					Status:            1,
					RepaymentSchedule: "weekly",
					Tenor:             2,
					RepaymentAmount:   1000,
					CreatedAt:         time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:         time.Time{},
				}, nil)
				f.DBRepo.EXPECT().SelectRepaymentCountByLoanId(gomock.Any(), int64(1)).Return(int(0), errs.NewWithMessage(http.StatusNotFound, ""))
				time := time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC)
				f.Clock.EXPECT().Now().Return(time).Times(2)

			},
			want: &entities.RepaymentInquiry{
				LoanId:          1,
				LoanReferenceId: "",
				LoanStatus:      "active",
				RepaymentNeeded: []entities.RepaymentNeeded{
					{
						Amount:  1000,
						DueDate: time.Date(2000, time.December, 8, 0, 0, 0, 0, time.UTC),
						IsLate:  false,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "status still active",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "reference",
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            1000,
					RatePercentage:    0,
					Status:            1,
					RepaymentSchedule: "weekly",
					Tenor:             2,
					RepaymentAmount:   1000,
					CreatedAt:         time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:         time.Time{},
				}, nil)
				f.DBRepo.EXPECT().SelectRepaymentCountByLoanId(gomock.Any(), int64(1)).Return(int(1), nil)
				time := time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC)
				f.Clock.EXPECT().Now().Return(time).Times(2)

			},
			want: &entities.RepaymentInquiry{
				LoanId:          1,
				LoanReferenceId: "",
				LoanStatus:      "active",
				RepaymentNeeded: []entities.RepaymentNeeded{
					{
						Amount:  1000,
						DueDate: time.Date(2000, time.December, 15, 0, 0, 0, 0, time.UTC),
						IsLate:  false,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error parameter",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "",
			},
			mock: func(f fields, args input) {
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error select loan",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "reference",
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param).Return(nil, errs.NewWithMessage(http.StatusInternalServerError, ""))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error select count",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: "reference",
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            1000,
					RatePercentage:    0,
					Status:            1,
					RepaymentSchedule: "weekly",
					Tenor:             2,
					RepaymentAmount:   1000,
					CreatedAt:         time.Date(2000, 12, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:         time.Time{},
				}, nil)
				f.DBRepo.EXPECT().SelectRepaymentCountByLoanId(gomock.Any(), int64(1)).Return(int(1), errs.NewWithMessage(http.StatusInternalServerError, ""))

			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			u := BillingUseCase{
				DBRepo: f.DBRepo,
				Clock:  f.Clock,
			}
			tt.mock(f, tt.input)

			got, err := u.GetRepaymentInquiryByLoanReferenceId(tt.input.ctx, tt.input.param)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestBillingUseCase_MakePayment(t *testing.T) {
	type input struct {
		ctx   context.Context
		param entities.RepaymentRequest
	}
	type fields struct {
		DBRepo *mock_usecase.MockDBRepository
		Clock  *mock_domain.MockClock
	}
	tests := []struct {
		name    string
		fields  func(ctrl *gomock.Controller) fields
		input   input
		mock    func(ctrl *gomock.Controller, f fields, input input)
		want    int64
		wantErr bool
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.RepaymentRequest{
					LoanReferenceId:      "reference",
					RepaymentReferenceId: "repaymentReference",
					Amount:               1000,
				},
			},
			mock: func(ctrl *gomock.Controller, f fields, args input) {
				f.DBRepo.EXPECT().SelectRepaymentByReferenceId(gomock.Any(), args.param.RepaymentReferenceId).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.LoanReferenceId).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            2000,
					RatePercentage:    0,
					Status:            entities.LoanStatusActive,
					RepaymentSchedule: "weekly",
					Tenor:             2,
					RepaymentAmount:   1000,
					CreatedAt:         time.Time{},
					UpdatedAt:         time.Time{},
				}, nil)
				f.DBRepo.EXPECT().SelectTotalRepaymentAmountByLoanId(gomock.Any(), int64(1)).Return(int64(0), errs.NewWithMessage(http.StatusNotFound, ""))
				tx := mock_domain.NewMockAtomicTransaction(ctrl)
				tx.EXPECT().Commit().Return(nil)
				tx.EXPECT().Rollback().Return(nil)
				f.DBRepo.EXPECT().BeginTx(gomock.Any()).Return(tx, nil)
				f.DBRepo.EXPECT().CreateRepayment(gomock.Any(), gomock.Any(), entities.Repayment{
					LoanId:      1,
					ReferenceId: "repaymentReference",
					Amount:      1000,
				}).Return(int64(1), nil)

			},
			want:    1,
			wantErr: false,
		},
		{
			name: "error parameter",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: entities.RepaymentRequest{},
			},
			mock: func(ctrl *gomock.Controller, f fields, args input) {
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error already exist",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.RepaymentRequest{
					LoanReferenceId:      "reference",
					RepaymentReferenceId: "repaymentReference",
					Amount:               100,
				},
			},
			mock: func(ctrl *gomock.Controller, f fields, args input) {
				f.DBRepo.EXPECT().SelectRepaymentByReferenceId(gomock.Any(), args.param.RepaymentReferenceId).Return(nil, nil)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error select repayment",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.RepaymentRequest{
					LoanReferenceId:      "reference",
					RepaymentReferenceId: "repaymentReference",
					Amount:               100,
				},
			},
			mock: func(ctrl *gomock.Controller, f fields, args input) {
				f.DBRepo.EXPECT().SelectRepaymentByReferenceId(gomock.Any(), args.param.RepaymentReferenceId).Return(nil, errs.NewWithMessage(http.StatusInternalServerError, ""))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error select loan",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.RepaymentRequest{
					LoanReferenceId:      "reference",
					RepaymentReferenceId: "repaymentReference",
					Amount:               100,
				},
			},
			mock: func(ctrl *gomock.Controller, f fields, args input) {
				f.DBRepo.EXPECT().SelectRepaymentByReferenceId(gomock.Any(), args.param.RepaymentReferenceId).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.LoanReferenceId).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            2000,
					RatePercentage:    0,
					Status:            entities.LoanStatusActive,
					RepaymentSchedule: "weekly",
					Tenor:             2,
					RepaymentAmount:   1000,
					CreatedAt:         time.Time{},
					UpdatedAt:         time.Time{},
				}, errs.NewWithMessage(http.StatusNotFound, ""))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error select repayment",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.RepaymentRequest{
					LoanReferenceId:      "reference",
					RepaymentReferenceId: "repaymentReference",
					Amount:               100,
				},
			},
			mock: func(ctrl *gomock.Controller, f fields, args input) {
				f.DBRepo.EXPECT().SelectRepaymentByReferenceId(gomock.Any(), args.param.RepaymentReferenceId).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.LoanReferenceId).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            2000,
					RatePercentage:    0,
					Status:            entities.LoanStatusActive,
					RepaymentSchedule: "weekly",
					Tenor:             2,
					RepaymentAmount:   1000,
					CreatedAt:         time.Time{},
					UpdatedAt:         time.Time{},
				}, nil)
				f.DBRepo.EXPECT().SelectTotalRepaymentAmountByLoanId(gomock.Any(), int64(1)).Return(int64(0), errs.NewWithMessage(http.StatusInternalServerError, ""))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error loan status",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.RepaymentRequest{
					LoanReferenceId:      "reference",
					RepaymentReferenceId: "repaymentReference",
					Amount:               100,
				},
			},
			mock: func(ctrl *gomock.Controller, f fields, args input) {
				f.DBRepo.EXPECT().SelectRepaymentByReferenceId(gomock.Any(), args.param.RepaymentReferenceId).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.LoanReferenceId).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            2000,
					RatePercentage:    0,
					Status:            entities.LoanStatusCompleted,
					RepaymentSchedule: "weekly",
					Tenor:             2,
					RepaymentAmount:   1000,
					CreatedAt:         time.Time{},
					UpdatedAt:         time.Time{},
				}, nil)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error amount",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.RepaymentRequest{
					LoanReferenceId:      "reference",
					RepaymentReferenceId: "repaymentReference",
					Amount:               100,
				},
			},
			mock: func(ctrl *gomock.Controller, f fields, args input) {
				f.DBRepo.EXPECT().SelectRepaymentByReferenceId(gomock.Any(), args.param.RepaymentReferenceId).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.LoanReferenceId).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            2000,
					RatePercentage:    0,
					Status:            entities.LoanStatusActive,
					RepaymentSchedule: "weekly",
					Tenor:             2,
					RepaymentAmount:   1000,
					CreatedAt:         time.Time{},
					UpdatedAt:         time.Time{},
				}, nil)
				f.DBRepo.EXPECT().SelectTotalRepaymentAmountByLoanId(gomock.Any(), int64(1)).Return(int64(0), errs.NewWithMessage(http.StatusNotFound, ""))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error begin",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.RepaymentRequest{
					LoanReferenceId:      "reference",
					RepaymentReferenceId: "repaymentReference",
					Amount:               1000,
				},
			},
			mock: func(ctrl *gomock.Controller, f fields, args input) {
				f.DBRepo.EXPECT().SelectRepaymentByReferenceId(gomock.Any(), args.param.RepaymentReferenceId).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.LoanReferenceId).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            2000,
					RatePercentage:    0,
					Status:            entities.LoanStatusActive,
					RepaymentSchedule: "weekly",
					Tenor:             2,
					RepaymentAmount:   1000,
					CreatedAt:         time.Time{},
					UpdatedAt:         time.Time{},
				}, nil)
				f.DBRepo.EXPECT().SelectTotalRepaymentAmountByLoanId(gomock.Any(), int64(1)).Return(int64(0), errs.NewWithMessage(http.StatusNotFound, ""))
				tx := mock_domain.NewMockAtomicTransaction(ctrl)
				f.DBRepo.EXPECT().BeginTx(gomock.Any()).Return(tx, errs.NewWithMessage(http.StatusInternalServerError, ""))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error create",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.RepaymentRequest{
					LoanReferenceId:      "reference",
					RepaymentReferenceId: "repaymentReference",
					Amount:               1000,
				},
			},
			mock: func(ctrl *gomock.Controller, f fields, args input) {
				f.DBRepo.EXPECT().SelectRepaymentByReferenceId(gomock.Any(), args.param.RepaymentReferenceId).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.LoanReferenceId).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            2000,
					RatePercentage:    0,
					Status:            entities.LoanStatusActive,
					RepaymentSchedule: "weekly",
					Tenor:             2,
					RepaymentAmount:   1000,
					CreatedAt:         time.Time{},
					UpdatedAt:         time.Time{},
				}, nil)
				f.DBRepo.EXPECT().SelectTotalRepaymentAmountByLoanId(gomock.Any(), int64(1)).Return(int64(0), errs.NewWithMessage(http.StatusNotFound, ""))
				tx := mock_domain.NewMockAtomicTransaction(ctrl)
				tx.EXPECT().Rollback().Return(nil)

				f.DBRepo.EXPECT().BeginTx(gomock.Any()).Return(tx, nil)
				f.DBRepo.EXPECT().CreateRepayment(gomock.Any(), gomock.Any(), entities.Repayment{
					LoanId:      1,
					ReferenceId: "repaymentReference",
					Amount:      1000,
				}).Return(int64(1), errs.NewWithMessage(http.StatusNotFound, ""))

			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error commit",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx: context.Background(),
				param: entities.RepaymentRequest{
					LoanReferenceId:      "reference",
					RepaymentReferenceId: "repaymentReference",
					Amount:               1000,
				},
			},
			mock: func(ctrl *gomock.Controller, f fields, args input) {
				f.DBRepo.EXPECT().SelectRepaymentByReferenceId(gomock.Any(), args.param.RepaymentReferenceId).Return(nil, errs.NewWithMessage(http.StatusNotFound, ""))
				f.DBRepo.EXPECT().SelectLoanByReferenceId(gomock.Any(), args.param.LoanReferenceId).Return(&entities.Loan{
					Id:                1,
					ReferenceId:       "",
					UserId:            0,
					Amount:            2000,
					RatePercentage:    0,
					Status:            entities.LoanStatusActive,
					RepaymentSchedule: "weekly",
					Tenor:             2,
					RepaymentAmount:   1000,
					CreatedAt:         time.Time{},
					UpdatedAt:         time.Time{},
				}, nil)
				f.DBRepo.EXPECT().SelectTotalRepaymentAmountByLoanId(gomock.Any(), int64(1)).Return(int64(0), errs.NewWithMessage(http.StatusNotFound, ""))
				tx := mock_domain.NewMockAtomicTransaction(ctrl)
				tx.EXPECT().Commit().Return(errs.NewWithMessage(http.StatusInternalServerError, ""))
				tx.EXPECT().Rollback().Return(nil)

				f.DBRepo.EXPECT().BeginTx(gomock.Any()).Return(tx, nil)
				f.DBRepo.EXPECT().CreateRepayment(gomock.Any(), gomock.Any(), entities.Repayment{
					LoanId:      1,
					ReferenceId: "repaymentReference",
					Amount:      1000,
				}).Return(int64(1), nil)

			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			u := BillingUseCase{
				DBRepo: f.DBRepo,
				Clock:  f.Clock,
			}
			tt.mock(ctrl, f, tt.input)

			got, err := u.MakePayment(tt.input.ctx, tt.input.param)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestBillingUseCase_GetLoanListByUserId(t *testing.T) {
	type input struct {
		ctx   context.Context
		param int64
	}
	type fields struct {
		DBRepo *mock_usecase.MockDBRepository
		Clock  *mock_domain.MockClock
	}
	tests := []struct {
		name    string
		fields  func(ctrl *gomock.Controller) fields
		input   input
		mock    func(f fields, input input)
		want    *[]entities.Loan
		wantErr bool
	}{
		{
			name: "success",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: 1,
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByUserId(gomock.Any(), args.param).Return(&[]entities.Loan{}, nil)
			},
			want:    &[]entities.Loan{},
			wantErr: false,
		},
		{
			name: "error param",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: 0,
			},
			mock: func(f fields, args input) {
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error db",
			fields: func(ctrl *gomock.Controller) fields {
				return fields{
					DBRepo: mock_usecase.NewMockDBRepository(ctrl),
					Clock:  mock_domain.NewMockClock(ctrl),
				}
			},
			input: input{
				ctx:   context.Background(),
				param: 1,
			},
			mock: func(f fields, args input) {
				f.DBRepo.EXPECT().SelectLoanByUserId(gomock.Any(), args.param).Return(&[]entities.Loan{}, errs.NewWithMessage(http.StatusInternalServerError, ""))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := tt.fields(ctrl)
			u := BillingUseCase{
				DBRepo: f.DBRepo,
				Clock:  f.Clock,
			}
			tt.mock(f, tt.input)

			got, err := u.GetLoanListByUserId(tt.input.ctx, tt.input.param)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.Nil(t, err)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
