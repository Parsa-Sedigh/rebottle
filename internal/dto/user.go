package dto

import "time"

type User struct {
	ID             int
	Phone          string
	FirstName      string
	LastName       string
	Email          string
	Credit         uint16
	Status         string
	EmailStatus    string
	Province       string
	City           string
	Street         string
	Alley          string
	ApartmentPlate uint16
	ApartmentNo    uint16
	PostalCode     string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateUser struct {
	Phone          string `json:"phone" validate:"required,min=11,max=11,phone"`
	FirstName      string `json:"first_name" validate:"required,min=3"`
	LastName       string `json:"last_name" validate:"required,min=3"`
	Email          string `json:"email,omitempty" validate:"omitempty,email"`
	Password       string `json:"password" validate:"min=6"`
	Province       string `json:"province" validate:"required"`
	City           string `json:"city" validate:"required"`
	Street         string `json:"street" validate:"required"`
	Alley          string `json:"alley,omitempty"`
	ApartmentPlate int    `json:"apartment_plate,omitempty" validate:"required"`
	ApartmentNo    int    `json:"apartment_no,omitempty" validate:"required"`
	PostalCode     string `json:"postal_code" validate:"required"`
}

type UpdateUser struct {
	FirstName      string `json:"first_name,omitempty" validate:"min=3"`
	LastName       string `json:"last_name,omitempty" validate:"min=3"`
	Email          string `json:"email,omitempty" validate:"omitempty,email"`
	Province       string `json:"province,omitempty"`
	City           string `json:"city,omitempty"`
	Street         string `json:"street,omitempty"`
	Alley          string `json:"alley,omitempty"`
	ApartmentPlate uint16 `json:"apartment_plate,omitempty"`
	ApartmentNo    uint16 `json:"apartment_no,omitempty"`
	PostalCode     string `json:"postal_code,omitempty"`
}

type GetUserResponse struct {
	ID             int    `json:"id"`
	Phone          string `json:"phone"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Credit         uint16 `json:"credit"`
	Status         string `json:"status"`
	EmailStatus    string `json:"email_status"`
	Province       string `json:"province"`
	City           string `json:"city"`
	Street         string `json:"street"`
	Alley          string `json:"alley"`
	ApartmentPlate uint16 `json:"apartment_plate"`
	ApartmentNo    uint16 `json:"apartment_no"`
	PostalCode     string `json:"postal_code"`
}
