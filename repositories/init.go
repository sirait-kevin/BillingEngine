package repositories

import "database/sql"

type UserRepository struct {
	DB *sql.DB
}
