package dto

type UpdateDriverRequest struct {
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
