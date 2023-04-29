package otp

import "math/rand"

var letterRunes = []rune("0123456789")

func GenerateOTP(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
