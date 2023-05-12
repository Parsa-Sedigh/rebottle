package dto

import "time"

type CreatePickupRequest struct {
	UserID int     `json:"user_id"`
	Time   int64   `json:"time"`
	Weight float64 `json:"weight"`
	Note   string  `json:"note"`
}

type CreatePickupRequestValidation struct {
	UserID int       `validate:"required"`
	Time   time.Time `validate:"required,gt"`
	Weight float64   `validate:"required"`
	Note   string
}

type UpdatePickupRequest struct {
	ID     int     `json:"id"`
	Time   int64   `json:"time"`
	Weight float32 `json:"weight"`
	Note   string  `json:"note"`
}

type UpdatePickupRequestValidation struct {
	ID     int     `validate:"required"`
	Time   int64   `validate:"required,gt"`
	Weight float32 `validate:"required"`
	Note   string
}

type CancelPickupRequest struct {
	ID int `json:"id" validate:"required"`
}
