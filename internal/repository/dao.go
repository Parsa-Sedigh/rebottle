package repository

import (
	"database/sql"
	"github.com/alexedwards/scs/v2"
	"time"
)

const dbTimeout = 3 * time.Second

type DAO interface {
	NewUserRepository() UserRepository
	NewAuthRepository() AuthRepository
	NewDriverRepository() DriverRepository
}

type dao struct {
	DB      *sql.DB
	session *scs.SessionManager
}

func NewDAO(DB *sql.DB) DAO {
	return &dao{
		DB: DB,
	}
}

func (d *dao) NewUserRepository() UserRepository {
	return &userRepository{DB: d.DB}
}

func (d *dao) NewAuthRepository() AuthRepository {
	return &authRepository{DB: d.DB}
}

func (d *dao) NewDriverRepository() DriverRepository {
	return &driverRepository{DB: d.DB}
}
