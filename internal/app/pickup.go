package app

import (
	"errors"
	"github.com/Parsa-Sedigh/rebottle/internal/appjwt"
	"github.com/Parsa-Sedigh/rebottle/internal/dto"
	"github.com/Parsa-Sedigh/rebottle/internal/models"
	"github.com/Parsa-Sedigh/rebottle/pkg/jsonutil"
	"github.com/Parsa-Sedigh/rebottle/pkg/validation"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

// TODO: Filters
func (app *application) GetPickups(w http.ResponseWriter, r *http.Request) {
	pickups, err := app.DB.GetUserPickups(int(r.Context().Value("JWTData").(appjwt.JWTData).UserID))
	if err != nil {
		jsonutil.BadRequest(w, r, err)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, pickups)
}

func (app *application) GetPickup(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	if ID == "" {
		jsonutil.BadRequest(w, r, errors.New("specify a pickup id"))
		return
	}

	IDNum, err := strconv.Atoi(ID)
	if err != nil {
		jsonutil.BadRequest(w, r, err)
		return
	}

	p, err := app.DB.GetPickup(IDNum, int(r.Context().Value("JWTData").(appjwt.JWTData).UserID))
	if err != nil {
		jsonutil.BadRequest(w, r, err)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, p)
}

func (app *application) CreatePickup(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreatePickupRequest

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		jsonutil.BadRequest(w, r, err)
		return
	}

	if translatedErr := validation.ValidatePayload(app.Validate, app.Translator, dto.CreatePickupRequestValidation{
		UserID: payload.UserID,
		Time:   time.UnixMilli(payload.Time),
		Weight: payload.Weight,
		Note:   payload.Note,
	}); translatedErr != nil {
		jsonutil.WriteJSON(w, http.StatusBadRequest, jsonutil.Resp{
			Error:   true,
			Message: "Some of the fields have error",
			Data:    translatedErr,
		})
		return
	}

	var response jsonutil.Resp

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
		jsonutil.WriteJSON(w, http.StatusInternalServerError, response)
		return
	}

	response.Error = false
	response.Message = "Pickup created"
	response.Data = p

	err = jsonutil.WriteJSON(w, http.StatusCreated, response)
	if err != nil {
		response.Error = true
		response.Message = "Internal server error"
		jsonutil.WriteJSON(w, http.StatusInternalServerError, response)
	}
}

func (app *application) UpdatePickup(w http.ResponseWriter, r *http.Request) {
	var payload dto.UpdatePickupRequest

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		jsonutil.BadRequest(w, r, err)
		return
	}

	if translatedErr := validation.ValidatePayload(app.Validate, app.Translator, dto.UpdatePickupRequestValidation{
		ID:     payload.ID,
		Time:   payload.Time,
		Weight: payload.Weight,
		Note:   payload.Note,
	}); translatedErr != nil {
		jsonutil.WriteJSON(w, http.StatusBadRequest, jsonutil.Resp{
			Error:   true,
			Message: "Some of the fields have error",
			Data:    translatedErr,
		})
		return
	}

	p, err := app.DB.GetPickup(payload.ID, int(r.Context().Value("JWTData").(appjwt.JWTData).UserID))
	if err != nil {
		jsonutil.BadRequest(w, r, err)
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
		jsonutil.BadRequest(w, r, err)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, updatedPickup)
}

func (app *application) CancelPickup(w http.ResponseWriter, r *http.Request) {
	var payload dto.CancelPickupRequest

	err := jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		jsonutil.BadRequest(w, r, err)
		return
	}

	userID := int(r.Context().Value("JWTData").(appjwt.JWTData).UserID)

	p, err := app.DB.GetPickup(payload.ID, userID)
	if err != nil {
		jsonutil.BadRequest(w, r, err)
		return
	}

	cancellationStatus := "pickup is already cancelled"
	if p.Status == models.StatusPickupCancelledByUser {
		cancellationStatus += " by user"
	} else if p.Status == models.StatusPickupCancelledBySystem {
		cancellationStatus += " by system"
	}

	if p.Status == models.StatusPickupCancelledByUser || p.Status == models.StatusPickupCancelledBySystem {
		jsonutil.BadRequest(w, r, errors.New(cancellationStatus))
		return
	}

	err = app.DB.CancelPickup(p.ID, userID, true)
	if err != nil {
		jsonutil.BadRequest(w, r, err)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Resp{
		Error:   false,
		Message: "cancelled pickup by user successfully",
	})
}
