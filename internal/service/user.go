package service

import (
	"github.com/Parsa-Sedigh/rebottle/internal/dto"
	"github.com/Parsa-Sedigh/rebottle/internal/model"
	"github.com/Parsa-Sedigh/rebottle/internal/repository"
)

type UserService interface {
	GetUserByID(id int) (dto.User, error)
	GetUserByPhone(phone string) (dto.User, error)
	CreateUser(user dto.CreateUser) (int, error)
	UpdateUser(user dto.UpdateUser) (dto.User, error)
}

type userService struct {
	dao repository.DAO
}

func NewUserService(dao repository.DAO) UserService {
	return &userService{dao: dao}
}

func UserModelToDTO(user model.User) dto.User {
	return dto.User{
		ID:             user.ID,
		Phone:          user.Phone,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Credit:         user.Credit,
		Status:         user.Status,
		EmailStatus:    user.EmailStatus,
		Province:       user.Province,
		City:           user.City,
		Street:         user.Street,
		Alley:          user.Alley,
		ApartmentPlate: user.ApartmentPlate,
		ApartmentNo:    user.ApartmentNo,
		PostalCode:     user.PostalCode,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func (u *userService) GetUserByID(id int) (dto.User, error) {
	user, err := u.dao.NewUserRepository().GetUserByID(id)
	if err != nil {
		return dto.User{}, err
	}

	return UserModelToDTO(user), nil
}

func (u *userService) GetUserByPhone(phone string) (dto.User, error) {
	user, err := u.dao.NewUserRepository().GetUserByPhone(phone)
	if err != nil {
		return dto.User{}, err
	}

	// TODO: omit password field from model.user here
	return UserModelToDTO(user), nil
}

func (u *userService) CreateUser(user dto.CreateUser) (int, error) {
	id, err := u.dao.NewUserRepository().CreateUser(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *userService) UpdateUser(user dto.UpdateUser) (dto.User, error) {
	return dto.User{}, nil
}
