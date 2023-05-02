package main

import (
	"database/sql"
	"errors"
	"github.com/Parsa-Sedigh/rebottle/internal/appjwt"
	"github.com/Parsa-Sedigh/rebottle/internal/models"
	"github.com/Parsa-Sedigh/rebottle/internal/otp"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CreatePickupRequest struct {
	UserID int     `json:"user_id" validate:"required"`
	Time   float64 `json:"time" validate:"required"`
	Weight float64 `json:"weight" validate:"required"`
	Note   string  `json:"note"`
}

type UpdatePickupRequest struct {
	ID     int       `json:"id" validate:"required"`
	Time   time.Time `json:"time" validate:"required"`
	Weight float32   `json:"weight" validate:"required"`
	Note   string    `json:"note"`
}

type VerifyUserSignupRequest struct {
	Phone string `json:"phone" validate:"required,min=11,max=11,phone"`
	OTP   string `json:"otp" validate:"required,min=6,max=6"`
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
	Email string `json:"email" validate:"required,email"`
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

	var response Resp

	p, err := app.DB.InsertPickup(models.Pickup{
		UserID: payload.UserID,
		Time:   time.UnixMilli(int64(payload.Time)),
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

	p, err := app.DB.GetPickup(payload.ID, int(r.Context().Value("JWTData").(appjwt.JWTData).UserID))
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	updatedPickup, err := app.DB.UpdatePickup(models.UpdatePickupParams{
		ID:     payload.ID,
		UserID: p.UserID,
		Time:   payload.Time,
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

	err = app.Validate.Struct(payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// check
	user, err := app.DB.GetUserByPhone(payload.Phone)
	if err != nil && err != sql.ErrNoRows {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// if user exists and has a status other than inactive:
	if user.ID > 0 && user.Status != "inactive" {
		app.errorJSON(w, errors.New("user already exists"), http.StatusBadRequest)
		return
	}

	// insert user, if not already inserted(with default inactive status)
	if err == sql.ErrNoRows {
		_, err = app.DB.InsertUser(payload)
		if err != nil {
			app.errorJSON(w, err, http.StatusBadRequest)
			return
		}
	}

	app.Session.Put(r.Context(), "otp", otp.GenerateOTP(6))

	// TODO: send the validation email and SMS, so that user can verify both, but the SMS verification is necessary for the user to be registered

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

	u, err := app.DB.GetUserByPhone(payload.Phone)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// TODO: Generate JWT and send it back (http cookie or in response?)
	token, err := appjwt.Generate(u.ID)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	app.writeJSON(w, http.StatusOK, Resp{
		Error:   false,
		Message: "Authenticated",
		Data: struct {
			Token string `json:"token"`
		}{Token: token},
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
		Data:  u,
	})
}

func (app *application) SendResetPasswordOTP(w http.ResponseWriter, r *http.Request) {
	var payload SendResetPasswordOTPRequest

	app.readJSON(w, r, &payload)

	u, err := app.DB.GetUserByEmail(payload.Email)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest) // TODO: test badRequest method here
		return
	}

	if u.EmailStatus != models.StatusUserEmailActive {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// TODO: Set OTP for resetting password in session

	// TODO: send OTP SMS for resetting password
}

func (app *application) ResetPassword(w http.ResponseWriter, r *http.Request) {
	// TODO: Check OTP with the one in session and if it was correct, reset the password and update it in DB

}