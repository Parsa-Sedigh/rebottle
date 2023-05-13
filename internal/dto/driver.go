package dto

import "time"

type Driver struct {
	ID             int       `json:"id"`
	Phone          string    `json:"phone"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Status         string    `json:"status"`
	Email          string    `json:"email"`
	EmailStatus    string    `json:"email_status"`
	Province       string    `json:"province"`
	City           string    `json:"city"`
	Street         string    `json:"street"`
	Alley          string    `json:"alley"`
	ApartmentPlate uint16    `json:"apartment_plate"`
	ApartmentNo    uint16    `json:"apartment_no"`
	PostalCode     string    `json:"postal_code"`
	LicenseNo      string    `json:"license_no"` // TODO: better validation for license number in Iran?
	LicenseStatus  string    `json:"license_status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UpdateDriver struct {
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	LicenseNo      string `json:"license_no" validate:"required"` // TODO: better validation for license number in Iran?
	Province       string `json:"province" validate:"required"`
	City           string `json:"city" validate:"required"`
	Street         string `json:"street" validate:"required"`
	Alley          string `json:"alley"`
	ApartmentPlate uint16 `json:"apartment_plate" validate:"required"`
	ApartmentNo    uint16 `json:"apartment_no" validate:"required"`
	PostalCode     string `json:"postal_code" validate:"required"`
}

type CreateDriver struct {
	Phone          string `json:"phone" validate:"required,min=11,max=11,phone"`
	Password       string `json:"password" validate:"required,min=6"`
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	LicenseNo      string `json:"license_no" validate:"required"` // TODO: better validation for license number in Iran?
	Province       string `json:"province" validate:"required"`
	City           string `json:"city" validate:"required"`
	Street         string `json:"street" validate:"required"`
	Alley          string `json:"alley"`
	ApartmentPlate uint16 `json:"apartment_plate" validate:"required"`
	ApartmentNo    uint16 `json:"apartment_no" validate:"required"`
	PostalCode     string `json:"postal_code" validate:"required"`
}
