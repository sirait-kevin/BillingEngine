package repositories

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DBRepository struct {
	DB *sqlx.DB
}

func (r *DBRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.DB.BeginTx(ctx, nil)
}
