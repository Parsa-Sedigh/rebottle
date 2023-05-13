package app

import (
	"errors"
	"fmt"
	"github.com/Parsa-Sedigh/rebottle/internal/appjwt"
	"github.com/Parsa-Sedigh/rebottle/internal/dto"
	"github.com/Parsa-Sedigh/rebottle/internal/models"
	"github.com/Parsa-Sedigh/rebottle/internal/otp"
	"github.com/Parsa-Sedigh/rebottle/pkg/jsonutil"
	"github.com/Parsa-Sedigh/rebottle/pkg/serviceerr"
	"github.com/Parsa-Sedigh/rebottle/pkg/validation"
	"net/http"
	"strings"
)

func (app *application) SignupUser(w http.ResponseWriter, r *http.Request) {
	var payload models.SignupUserRequest
	fmt.Println("hello: ", app.userService, app.authService)

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
		return
	}

	/* TODO: Should be in service layer, but how we should pass app.Validate and app.Translator there? Should we pass them
	   when we instantiate the AuthService ?*/
	validation.ValidatePayload(app.Validate, app.Translator, payload)

	// check
	_, err = app.authService.SignupUser(r.Context(), dto.CreateUser{
		Phone:          payload.Phone,
		FirstName:      payload.FirstName,
		LastName:       payload.LastName,
		Email:          payload.Email,
		Province:       payload.Province,
		City:           payload.City,
		Street:         payload.Street,
		Alley:          payload.Alley,
		ApartmentPlate: payload.ApartmentPlate,
		ApartmentNo:    payload.ApartmentNo,
		PostalCode:     payload.PostalCode,
	})
	if err, ok := err.(serviceerr.ServiceErr); ok {
		jsonutil.ErrorJSON(w, app.logger, err, err.Status)
		return
	}

	resp := jsonutil.Resp{
		Error:   false,
		Message: "User created",
	}

	if err = jsonutil.WriteJSON(w, http.StatusCreated, resp); err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
	}
}

func (app *application) VerifyUserSignup(w http.ResponseWriter, r *http.Request) {
	var payload dto.VerifyUserSignupRequest

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
		return
	}

	err = app.authService.VerifyUserSignup(r.Context(), payload)
	if err, ok := err.(serviceerr.ServiceErr); ok {
		jsonutil.ErrorJSON(w, app.logger, err, err.Status)
		return
	}

	resp := jsonutil.Resp{
		Error:   false,
		Message: "User activated",
	}
	jsonutil.WriteJSON(w, http.StatusOK, resp)
}

func (app *application) authenticateToken(r *http.Request) (*models.User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	// one space occurs between Bearer and the token itself
	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no authorization header received")
	}

	token := headerParts[1]
	if len(token) != 26 {
		return nil, errors.New("authentication token wrong size")
	}

	// get user from JWT
	// get the user from the tokens table
	//user, err := app.DB.GetUserForToken(token)
	//if err != nil {
	//	return nil, errors.New("no matching user found")
	//}

	//return user, nil
	return &models.User{}, nil
}

// VerifyUserEmail TODO
func (app *application) VerifyUserEmail(w http.ResponseWriter, r *http.Request) {
	var payload dto.VerifyUserEmailRequest

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
		return
	}
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	var payload dto.LoginRequest

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
		return
	}

	if translatedErr := validation.ValidatePayload(app.Validate, app.Translator, payload); translatedErr != nil {
		jsonutil.WriteJSON(w, http.StatusBadRequest, jsonutil.Resp{
			Error:   true,
			Message: "Some of the fields have error",
			Data:    translatedErr,
		})
		return
	}

	//var u models.User
	//var d models.Driver
	//
	//if payload.IsDriver {
	//	d, err = app.DB.GetDriverByPhone(payload.Phone)
	//} else {
	//	u, err = app.DB.GetUserByPhone(payload.Phone)
	//}
	//
	//if err != nil {
	//	jsonutil.ErrorJSON(w, app.logger, errors.New("invalid credentials"), http.StatusBadRequest)
	//	return
	//}
	//
	//var isPasswordCorrect bool
	//if payload.IsDriver {
	//	isPasswordCorrect = password.CheckPasswordHash(payload.Password, d.Password)
	//} else {
	//	isPasswordCorrect = password.CheckPasswordHash(payload.Password, u.Password)
	//}
	//
	//if !isPasswordCorrect {
	//	jsonutil.ErrorJSON(w, app.logger, errors.New("invalid credentials"), http.StatusBadRequest)
	//	return
	//}
	//
	//var accessToken, refreshToken string
	//if payload.IsDriver {
	//	// TODO: Generate JWT and send it back (http cookie or in response?)
	//	accessToken, refreshToken, err = appjwt.GenerateWithMoreClaims(u.ID, map[string]any{"isDriver": true})
	//} else {
	//	accessToken, refreshToken, err = appjwt.Generate(u.ID)
	//}

	accessToken, refreshToken, err := app.authService.Login(payload)
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Resp{
		Error:   false,
		Message: "Authenticated",
		Data: struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}{AccessToken: accessToken, RefreshToken: refreshToken},
	})
}

func (app *application) SignupDriver(w http.ResponseWriter, r *http.Request) {
	var payload dto.SignupDriverRequest

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		jsonutil.BadRequest(w, r, err)
		return
	}

	validation.ValidatePayload(app.Validate, app.Translator, payload)

	driver, err := app.authService.SignupDriver(r.Context(), payload)

	if err, ok := err.(serviceerr.ServiceErr); ok {
		jsonutil.ErrorJSON(w, app.logger, err, err.Status)
		return
	}

	err = jsonutil.WriteJSON(w, http.StatusCreated, jsonutil.Resp{
		Error:   false,
		Message: "Driver created",
		Data:    driver,
	})
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
		return
	}
}

func (app *application) SendResetPasswordOTP(w http.ResponseWriter, r *http.Request) {
	var payload dto.SendResetPasswordOTPRequest

	jsonutil.ReadJSON(w, r, &payload)

	validation.ValidatePayload(app.Validate, app.Translator, payload)

	u, err := app.DB.GetUserByPhone(payload.Phone)
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest) // TODO: test badRequest method here
		return
	}

	if u.Status != models.StatusUserActive {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
		return
	}

	// TODO: Set OTP for resetting password in session
	resetPasswordOTP := otp.GenerateOTPCode(6)
	fmt.Println("signup otp: ", resetPasswordOTP)
	app.Session.Put(r.Context(), "resetPassword", dto.PhoneWithOTP{
		Phone: payload.Phone,
		OTP:   resetPasswordOTP,
	})

	// TODO: send OTP SMS for resetting password
}

func (app *application) VerifyResetPasswordOTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Check OTP with the one in session and if it was correct, reset the password and update it in DB
	// TODO: Should sending the OTP of resetting password be a separate step
	var payload dto.VerifyResetPasswordOTPRequest

	jsonutil.ReadJSON(w, r, &payload)

	validation.ValidatePayload(app.Validate, app.Translator, payload)

	if payload.OTP != app.Session.Get(r.Context(), "resetPassword").(dto.PhoneWithOTP).OTP {
		jsonutil.ErrorJSON(w, app.logger, errors.New("invalid OTP"), http.StatusBadRequest)
		return
	}

	// generate a token for sending it in url and pass it back in CompleteResetPassword
}

func (app *application) CompleteResetPassword(w http.ResponseWriter, r *http.Request) {
	var payload dto.CompleteResetPasswordRequest

	jsonutil.ReadJSON(w, r, &payload)

	validation.ValidatePayload(app.Validate, app.Translator, payload)

	// todo: Verify the token sent in query param(what should be in the claims? Maybe userID because we need it for updating the password)

	//resetPasswordSession := app.Session.Get(r.Context(), "resetPassword")
	//hashedPassword, err := password.HashPassword(payload.Password)
	//if err != nil {
	//	jsonutil.ErrorJSON(w, app.logger,err, http.StatusBadRequest)
	//	return
	//}

	//err := app.DB.UpdateUserPassword(hashedPassword)
	//if err != nil {
	//	jsonutil.ErrorJSON(w, app.logger,err, http.StatusBadRequest)
	//	return
	//}
}

// NewAuthTokens generates a new pair of access and refresh tokens
func (app *application) NewAuthTokens(w http.ResponseWriter, r *http.Request) {
	var payload dto.NewAuthTokensRequest

	jsonutil.ReadJSON(w, r, &payload)
	JWTData, err := appjwt.ExtractClaims(payload.RefreshToken)
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, errors.New("you're unauthorized"), http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := appjwt.Generate(int(JWTData.UserID))
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, errors.New("you're unauthorized"), http.StatusUnauthorized)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, dto.NewAuthTokensResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	})
}

func (app *application) VerifyDriverSignup(w http.ResponseWriter, r *http.Request) {
	var payload dto.VerifyDriverSignupRequest

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		jsonutil.BadRequest(w, r, err)
		return
	}

	OTPData, ok := app.Session.Get(r.Context(), "otpData").(dto.PhoneWithOTP)
	if !ok {
		jsonutil.ErrorJSON(w, app.logger, errors.New("please try again"), http.StatusInternalServerError)
		return
	}

	if payload.OTP != OTPData.OTP {
		jsonutil.ErrorJSON(w, app.logger, errors.New("invalid OTP"), http.StatusBadRequest)
		return
	}

	driver, err := app.DB.GetDriverByPhone(OTPData.Phone)
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
		return
	}

	driver, err = app.DB.UpdateDriverStatus("active", driver.ID)
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := appjwt.GenerateWithMoreClaims(driver.ID, map[string]any{
		"isDriver": true,
	})
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, dto.NewAuthTokensResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	})
}
