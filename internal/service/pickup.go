package service

import (
	"github.com/Parsa-Sedigh/rebottle/internal/dto"
	"github.com/Parsa-Sedigh/rebottle/internal/model"
	"github.com/Parsa-Sedigh/rebottle/internal/repository"
	"github.com/Parsa-Sedigh/rebottle/pkg/serviceerr"
	"net/http"
)

type PickupService interface {
	CreatePickup(pickup dto.Pickup) error
	GetPickups(userID int) ([]dto.Pickup, error)
	GetPickup(id int, userID int) (dto.Pickup, error)
}

type pickupService struct {
	dao repository.DAO
}

func NewPickupService(dao repository.DAO) PickupService {
	return &pickupService{dao: dao}
}

func pickupModelToDTO(p model.Pickup) dto.Pickup {
	return dto.Pickup{
		ID:        p.ID,
		TruckID:   p.TruckID,
		UserID:    p.UserID,
		Time:      p.Time,
		Weight:    p.Weight,
		Note:      p.Note,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (p *pickupService) CreatePickup(pickup dto.Pickup) error {
	_, err := p.dao.NewPickupRepository().CreatePickup(pickup)
	if err != nil {
		return err
	}

	return nil
}

func (p *pickupService) GetPickups(userID int) ([]dto.Pickup, error) {
	pickups, err := p.dao.NewPickupRepository().GetUserPickups(userID)
	if err != nil {
		return []dto.Pickup{}, serviceerr.NewServiceErr("Sth went wrong", http.StatusBadRequest)
	}

	var dtoPickups []dto.Pickup

	for _, pickup := range pickups {
		dtoPickups = append(dtoPickups, pickupModelToDTO(pickup))
	}

	return dtoPickups, nil
}

func (p *pickupService) GetPickup(id int, userID int) (dto.Pickup, error) {
	pickup, err := p.dao.NewPickupRepository().GetUserPickup(id, userID)
	if err != nil {
		return dto.Pickup{}, serviceerr.NewServiceErr("Sth went wrong", http.StatusBadRequest)
	}

	return pickupModelToDTO(pickup), nil
}
