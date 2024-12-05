package repositories

import "github.com/jmoiron/sqlx"

type UserRepository struct {
	DB *sqlx.DB
}
