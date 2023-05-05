package otp

import "math/rand"

type OTP struct {
	APIKey string
}

var letterRunes = []rune("0123456789")

func NewOTP(APIKey string) OTP {
	return OTP{
		APIKey: APIKey,
	}
}

func GenerateOTPCode(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (o *OTP) SendSMS() {}
