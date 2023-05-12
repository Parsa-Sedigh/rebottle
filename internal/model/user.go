package model

import "time"

type User struct {
	ID             int
	Phone          string
	FirstName      string
	LastName       string
	Email          string
	Password       string
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
