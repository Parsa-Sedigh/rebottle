package appjwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type JWTData struct {
	Exp    float64
	UserID float64
}

func Generate(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":    time.Now().Add(3 * time.Minute).Unix(),
		"userID": userID,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Parse(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return JWTData{}, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return token, err
	}

	return token, nil
}

func ExtractClaims(tokenString string) (JWTData, error) {
	token, err := Parse(tokenString)
	if err != nil {
		return JWTData{}, err
	}
	if !token.Valid {
		return JWTData{}, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return JWTData{
			Exp:    claims["exp"].(float64),
			UserID: claims["userID"].(float64),
		}, nil
	}

	return JWTData{}, errors.New("unable to extract claims")
}
