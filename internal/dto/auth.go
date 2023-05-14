package dto

// TODO: Create a PhoneWithOTP struct and swaggerui it where ever a phone and otp fields are used
type PhoneWithOTP struct {
	Phone string `json:"phone" validate:"required,min=11,max=11,phone"`
	OTP   string `json:"otp" validate:"required,min=6,max=6,numeric"`
}

type VerifyUserSignupRequest struct {
	OTP string `json:"otp" validate:"required,min=6,max=6,numeric"`
}

type VerifyUserEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Hash  string `json:"hash" validate:"required"`
}

type LoginRequest struct {
	Phone    string `json:"phone" validate:"required,min=11,max=11,phone"`
	Password string `json:"password" validate:"required,min=6"`
	IsDriver bool   `json:"is_driver"`
}

type SignupDriverRequest struct {
	Phone          string `json:"phone" validate:"required,min=11,max=11,phone"`
	FirstName      string `json:"first_name" validate:"required"`
	LastName       string `json:"last_name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6"`
	LicenseNo      string `json:"license_no" validate:"required"` // TODO: better validation for license number in Iran?
	Province       string `json:"province" validate:"required"`
	City           string `json:"city" validate:"required"`
	Street         string `json:"street" validate:"required"`
	Alley          string `json:"alley"`
	ApartmentPlate uint16 `json:"apartment_plate" validate:"required"`
	ApartmentNo    uint16 `json:"apartment_no" validate:"required"`
	PostalCode     string `json:"postal_code" validate:"required"`
}

type SendResetPasswordOTPRequest struct {
	Phone string `json:"phone" validate:"required,min=11,max=11,phone"`
}

type VerifyResetPasswordOTPRequest struct {
	OTP string `json:"otp" validate:"required,min=6,max=6,numeric"`
}

type CompleteResetPasswordRequest struct {
	Password string `json:"password" validate:"required,min=6"`
}

type VerifyDriverSignupRequest struct {
	OTP string `json:"otp" validate:"required,min=6,max=6,numeric"`
}

type NewAuthTokensRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type NewAuthTokensResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}
