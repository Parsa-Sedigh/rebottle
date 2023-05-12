package model

import "time"

type Driver struct {
	ID             int
	Phone          string
	FirstName      string
	LastName       string
	Email          string
	Password       string
	Status         string
	EmailStatus    string
	Province       string
	City           string
	Street         string
	Alley          string
	ApartmentPlate uint16
	ApartmentNo    uint16
	PostalCode     string
	//UserID    int
	LicenseNo     string
	LicenseStatus string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
