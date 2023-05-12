package app

import (
	"github.com/Parsa-Sedigh/rebottle/internal/appjwt"
	"github.com/Parsa-Sedigh/rebottle/internal/dto"
	"github.com/Parsa-Sedigh/rebottle/pkg/jsonutil"
	"net/http"
)

func (app *application) GetUser(w http.ResponseWriter, r *http.Request) {
	u, err := app.DB.GetUserByID(int(r.Context().Value("JWTData").(appjwt.JWTData).UserID))
	if err != nil {
		jsonutil.BadRequest(w, r, err)
		return
	}

	jsonutil.WriteJSON(w, http.StatusOK, jsonutil.Resp{
		Error: false,
		Data: dto.GetUserResponse{
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
