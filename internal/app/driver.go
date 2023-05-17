package app

import (
	"github.com/Parsa-Sedigh/rebottle/internal/dto"
	"github.com/Parsa-Sedigh/rebottle/internal/models"
	"github.com/Parsa-Sedigh/rebottle/pkg/jsonutil"
	"github.com/Parsa-Sedigh/rebottle/pkg/validation"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (app *application) UpdateDriver(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")
	IDNum, err := strconv.Atoi(ID)
	if err != nil {
		jsonutil.BadRequest(w, r, err)
		return
	}

	var payload dto.UpdateDriver

	err = jsonutil.ReadJSON(w, r, &payload)
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
		return
	}

	validation.ValidatePayload(app.Validate, app.Translator, payload)

	driver, err := app.DB.UpdateDriver(models.UpdateDriverData{
		FirstName:      payload.FirstName,
		LastName:       payload.LastName,
		Email:          payload.Email,
		LicenseNo:      payload.LicenseNo,
		Province:       payload.Province,
		City:           payload.City,
		Street:         payload.Street,
		Alley:          payload.Alley,
		ApartmentPlate: payload.ApartmentPlate,
		ApartmentNo:    payload.ApartmentNo,
		PostalCode:     payload.PostalCode,
		ID:             IDNum,
	})
	if err != nil {
		jsonutil.ErrorJSON(w, app.logger, err, http.StatusBadRequest)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, driver)
}

func (app *application) InactiveDriver(w http.ResponseWriter, r *http.Request) {

}
