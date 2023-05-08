package appjwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type JWTData struct {
	Exp    float64
	UserID float64
}

func getDefaultClaims(userID int, expMinutes time.Duration) jwt.MapClaims {
	return jwt.MapClaims{
		"exp":    time.Now().Add(expMinutes * time.Minute).Unix(),
		"userID": userID,
	}
}

func generateAuthTokens(accessTokenClaims, refreshTokenClaims jwt.MapClaims) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Sign and get the complete encoded token as a string using the secret
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func Generate(userID int) (string, string, error) {
	return generateAuthTokens(getDefaultClaims(userID, time.Duration(1)), getDefaultClaims(userID, 2))
}

func GenerateWithMoreClaims(userID int, data map[string]any) (string, string, error) {
	accessTokenClaims := getDefaultClaims(userID, 1)

	for k, v := range data {
		fmt.Println("hello: ", k, v)

		accessTokenClaims[k] = v
	}

	return generateAuthTokens(accessTokenClaims, getDefaultClaims(userID, 2))
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
