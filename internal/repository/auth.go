package repository

import "database/sql"

type AuthRepository interface {
}

type authRepository struct {
	DB *sql.DB
}
