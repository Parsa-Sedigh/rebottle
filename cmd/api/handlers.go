package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Parsa-Sedigh/rebottle/internal/appjwt"
	"github.com/Parsa-Sedigh/rebottle/internal/models"
	"github.com/Parsa-Sedigh/rebottle/internal/otp"
	"github.com/Parsa-Sedigh/rebottle/internal/password"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO: Create a PhoneWithOTP struct and embed it where ever a phone and otp fields are used
type PhoneWithOTP struct {
	Phone string `json:"phone" validate:"required,min=11,max=11,phone"`
	OTP   string `json:"otp" validate:"required,min=6,max=6,numeric"`
}

type CreatePickupRequest struct {
	UserID int     `json:"user_id"`
	Time   int64   `json:"time"`
	Weight float64 `json:"weight"`
	Note   string  `json:"note"`
}

type CreatePickupRequestValidation struct {
	UserID int       `validate:"required"`
	Time   time.Time `validate:"required,gt"`
	Weight float64   `validate:"required"`
	Note   string
}

type UpdatePickupRequest struct {
	ID     int     `json:"id"`
	Time   int64   `json:"time"`
	Weight float32 `json:"weight"`
	Note   string  `json:"note"`
}

type UpdatePickupRequestValidation struct {
	ID     int     `validate:"required"`
	Time   int64   `validate:"required,gt"`
	Weight float32 `validate:"required"`
	Note   string
}

type VerifyUserSignupRequest struct {
	Phone string `json:"phone" validate:"required,min=11,max=11,phone"`
	OTP   string `json:"otp" validate:"required,min=6,max=6,numeric"`
}

type VerifyUserEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Hash  string `json:"hash" validate:"required"`
}

type LoginRequest struct {
	Phone    string `json:"phone" validate:"required,min=11,max=11,phone"`
	Password string `json:"password" validate:"required,min=6"`
}

type CancelPickupRequest struct {
	ID int `json:"id" validate:"required"`
}

type SendResetPasswordOTPRequest struct {
	Phone string `json:"phone" validate:"required,min=11,max=11,phone"`
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

type NewAuthTokensRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type NewAuthTokensResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type VerifyResetPasswordOTPRequest struct {
	OTP string `json:"otp" validate:"required,min=6,max=6,numeric"`
}

type CompleteResetPasswordRequest struct {
	Password string `json:"password" validate:"required,min=6"`
}

// TODO: Filters
func (app *application) GetPickups(w http.ResponseWriter, r *http.Request) {
	pickups, err := app.DB.GetUserPickups(int(r.Context().Value("JWTData").(appjwt.JWTData).UserID))
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, pickups)
}

func (app *application) GetPickup(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	if ID == "" {
		app.badRequest(w, r, errors.New("specify a pickup id"))
		return
	}

	IDNum, err := strconv.Atoi(ID)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	p, err := app.DB.GetPickup(IDNum, int(r.Context().Value("JWTData").(appjwt.JWTData).UserID))
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, p)
}

func (app *application) CreatePickup(w http.ResponseWriter, r *http.Request) {
	var payload CreatePickupRequest

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.validatePayload(w, CreatePickupRequestValidation{
		UserID: payload.UserID,
		Time:   time.UnixMilli(payload.Time),
		Weight: payload.Weight,
		Note:   payload.Note,
	})

	var response Resp

	p, err := app.DB.InsertPickup(models.Pickup{
		UserID: payload.UserID,
		Time:   time.UnixMilli(payload.Time),
		Weight: float32(payload.Weight),
		Note:   payload.Note,
		Status: models.StatusPickupWaiting,
	})
	if err != nil {
		response.Error = true
		response.Message = "Internal server error"
		app.writeJSON(w, http.StatusInternalServerError, response)
		return
	}

	response.Error = false
	response.Message = "Pickup created"
	response.Data = p

	err = app.writeJSON(w, http.StatusCreated, response)
	if err != nil {
		response.Error = true
		response.Message = "Internal server error"
		app.writeJSON(w, http.StatusInternalServerError, response)
	}
}

func (app *application) UpdatePickup(w http.ResponseWriter, r *http.Request) {
	var payload UpdatePickupRequest

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.validatePayload(w, UpdatePickupRequestValidation{
		ID:     payload.ID,
		Time:   payload.Time,
		Weight: payload.Weight,
		Note:   payload.Note,
	})

	p, err := app.DB.GetPickup(payload.ID, int(r.Context().Value("JWTData").(appjwt.JWTData).UserID))
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	updatedPickup, err := app.DB.UpdatePickup(models.UpdatePickupParams{
		ID:     payload.ID,
		UserID: p.UserID,
		Time:   time.UnixMilli(payload.Time),
		Weight: payload.Weight,
		Note:   payload.Note,
		Status: p.Status,
	})
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, updatedPickup)
}

func (app *application) SignupUser(w http.ResponseWriter, r *http.Request) {
	var payload models.SignupUserRequest

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	app.validatePayload(w, payload)

	// check
	user, err := app.DB.GetUserByPhone(payload.Phone)
	if err != nil && err != sql.ErrNoRows {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// if user exists:
	if user.ID > 0 {
		app.errorJSON(w, errors.New("user already exists"), http.StatusBadRequest)
		return
	}

	hashedPassword, err := password.HashPassword(payload.Password)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// insert user, if not already inserted(with default inactive status)
	if err == sql.ErrNoRows {
		_, err = app.DB.InsertUser(models.SignupUserRequest{
			Phone:          payload.Phone,
			FirstName:      payload.FirstName,
			LastName:       payload.LastName,
			Email:          payload.Email,
			Password:       hashedPassword,
			Province:       payload.Province,
			City:           payload.City,
			Street:         payload.Street,
			Alley:          payload.Alley,
			ApartmentPlate: payload.ApartmentPlate,
			ApartmentNo:    payload.ApartmentNo,
			PostalCode:     payload.PostalCode,
		})
		if err != nil {
			app.errorJSON(w, err, http.StatusBadRequest)
			return
		}
	}

	signupOTP := otp.GenerateOTPCode(6)
	fmt.Println("signup otp: ", signupOTP)
	app.Session.Put(r.Context(), "otp", signupOTP)

	/* TODO: send the validation email(IF provided) and SMS, so that user can verify both, but the SMS verification is necessary for the user to be registered.
	We can use the message field of Resp type and make it to have fa and en.*/

	resp := Resp{
		Error:   false,
		Message: "User created",
	}

	err = app.writeJSON(w, http.StatusCreated, resp)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
}

func (app *application) VerifyUserSignup(w http.ResponseWriter, r *http.Request) {
	var payload VerifyUserSignupRequest

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// check if OTP is correct
	if payload.OTP != app.Session.Get(r.Context(), "otp") {
		app.errorJSON(w, errors.New("invalid OTP"), http.StatusBadRequest)
		return
	}

	user, err := app.DB.GetUserByPhone(payload.Phone)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// update user status to active
	err = app.DB.UpdateUserStatus("active", user.ID)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	resp := Resp{
		Error:   false,
		Message: "User activated",
	}
	app.writeJSON(w, http.StatusOK, resp)
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

func (app *application) VerifyUserEmail(w http.ResponseWriter, r *http.Request) {
	var payload VerifyUserEmailRequest

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	var payload LoginRequest

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	app.validatePayload(w, payload)

	u, err := app.DB.GetUserByPhone(payload.Phone)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	isPasswordCorrect := password.CheckPasswordHash(payload.Password, u.Password)

	if !isPasswordCorrect {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	// TODO: Generate JWT and send it back (http cookie or in response?)
	accessToken, refreshToken, err := appjwt.Generate(u.ID)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusOK, Resp{
		Error:   false,
		Message: "Authenticated",
		Data: struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}{AccessToken: accessToken, RefreshToken: refreshToken},
	})
}

func (app *application) CancelPickup(w http.ResponseWriter, r *http.Request) {
	var payload CancelPickupRequest

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	userID := int(r.Context().Value("JWTData").(appjwt.JWTData).UserID)

	p, err := app.DB.GetPickup(payload.ID, userID)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	cancellationStatus := "pickup is already cancelled"
	if p.Status == models.StatusPickupCancelledByUser {
		cancellationStatus += " by user"
	} else if p.Status == models.StatusPickupCancelledBySystem {
		cancellationStatus += " by system"
	}

	if p.Status == models.StatusPickupCancelledByUser || p.Status == models.StatusPickupCancelledBySystem {
		app.badRequest(w, r, errors.New(cancellationStatus))
		return
	}

	err = app.DB.CancelPickup(p.ID, userID, true)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, Resp{
		Error:   false,
		Message: "cancelled pickup by user successfully",
	})
}

func (app *application) GetUser(w http.ResponseWriter, r *http.Request) {
	u, err := app.DB.GetUserByID(int(r.Context().Value("JWTData").(appjwt.JWTData).UserID))
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, Resp{
		Error: false,
		Data: GetUserResponse{
			ID:             u.ID,
			Phone:          u.Phone,
			FirstName:      u.FirstName,
			LastName:       u.LastName,
			Email:          u.Email,
			Credit:         u.Credit,
			Status:         u.Status,
			EmailStatus:    u.EmailStatus,
			Province:       u.Province,
			City:           u.City,
			Street:         u.Street,
			Alley:          u.Alley,
			ApartmentPlate: u.ApartmentPlate,
			ApartmentNo:    u.ApartmentNo,
			PostalCode:     u.PostalCode,
		},
	})
}

func (app *application) SendResetPasswordOTP(w http.ResponseWriter, r *http.Request) {
	var payload SendResetPasswordOTPRequest

	app.readJSON(w, r, &payload)

	app.validatePayload(w, payload)

	u, err := app.DB.GetUserByPhone(payload.Phone)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest) // TODO: test badRequest method here
		return
	}

	if u.Status != models.StatusUserActive {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// TODO: Set OTP for resetting password in session
	resetPasswordOTP := otp.GenerateOTPCode(6)
	fmt.Println("signup otp: ", resetPasswordOTP)
	app.Session.Put(r.Context(), "resetPassword", PhoneWithOTP{
		Phone: payload.Phone,
		OTP:   resetPasswordOTP,
	})

	// TODO: send OTP SMS for resetting password
}

func (app *application) VerifyResetPasswordOTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Check OTP with the one in session and if it was correct, reset the password and update it in DB
	// TODO: Should sending the OTP of resetting password be a separate step
	var payload VerifyResetPasswordOTPRequest

	app.readJSON(w, r, &payload)

	app.validatePayload(w, payload)

	if payload.OTP != app.Session.Get(r.Context(), "resetPassword").(PhoneWithOTP).OTP {
		app.errorJSON(w, errors.New("invalid OTP"), http.StatusBadRequest)
		return
	}

	// generate a token for sending it in url and pass it back in CompleteResetPassword
}

func (app *application) CompleteResetPassword(w http.ResponseWriter, r *http.Request) {
	var payload CompleteResetPasswordRequest

	app.readJSON(w, r, &payload)

	app.validatePayload(w, payload)

	// todo: Verify the token sent in query param(what should be in the claims? Maybe userID because we need it for updating the password)

	//resetPasswordSession := app.Session.Get(r.Context(), "resetPassword")
	//hashedPassword, err := password.HashPassword(payload.Password)
	//if err != nil {
	//	app.errorJSON(w, err, http.StatusBadRequest)
	//	return
	//}

	//err := app.DB.UpdateUserPassword(hashedPassword)
	//if err != nil {
	//	app.errorJSON(w, err, http.StatusBadRequest)
	//	return
	//}
}

// NewAuthTokens generates a new pair of access and refresh tokens
func (app *application) NewAuthTokens(w http.ResponseWriter, r *http.Request) {
	var payload NewAuthTokensRequest

	app.readJSON(w, r, &payload)
	JWTData, err := appjwt.ExtractClaims(payload.RefreshToken)
	if err != nil {
		app.errorJSON(w, errors.New("you're unauthorized"), http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := appjwt.Generate(int(JWTData.UserID))
	if err != nil {
		app.errorJSON(w, errors.New("you're unauthorized"), http.StatusUnauthorized)
		return
	}

	app.writeJSON(w, http.StatusOK, NewAuthTokensResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	})
}
