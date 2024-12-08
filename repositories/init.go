package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/sirait-kevin/BillingEngine/domain/interfaces"
)

type DBRepository struct {
	DB *sqlx.DB
}

func (r *DBRepository) BeginTx(ctx context.Context) (interfaces.AtomicTransaction, error) {
	return r.DB.BeginTx(ctx, nil)
}
