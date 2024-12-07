package repositories

import "github.com/jmoiron/sqlx"

type DBRepository struct {
	DB *sqlx.DB
}
