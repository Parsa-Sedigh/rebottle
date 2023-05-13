package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Parsa-Sedigh/rebottle/internal/appjwt"
	"github.com/Parsa-Sedigh/rebottle/internal/dto"
	"github.com/Parsa-Sedigh/rebottle/internal/model"
	"github.com/Parsa-Sedigh/rebottle/internal/models"
	"github.com/Parsa-Sedigh/rebottle/internal/otp"
	"github.com/Parsa-Sedigh/rebottle/internal/password"
	"github.com/Parsa-Sedigh/rebottle/internal/repository"
	"github.com/Parsa-Sedigh/rebottle/pkg/serviceerr"
	"github.com/alexedwards/scs/v2"
	"net/http"
)

type AuthService interface {
	SignupUser(ctx context.Context, payload dto.CreateUser) (dto.User, error)
	SignupDriver(ctx context.Context, payload dto.SignupDriverRequest) (dto.Driver, error)
	VerifyUserSignup(ctx context.Context, payload dto.VerifyUserSignupRequest) error
	Login(payload dto.LoginRequest) (string, string, error)
}

type authService struct {
	dao     repository.DAO
	session *scs.SessionManager
}

func driverModelToDTO(d model.Driver) dto.Driver {
	return dto.Driver{
		ID:             d.ID,
		Phone:          d.Phone,
		FirstName:      d.FirstName,
		LastName:       d.LastName,
		Email:          d.Email,
		EmailStatus:    d.EmailStatus,
		Status:         d.Status,
		Province:       d.Province,
		City:           d.City,
		Street:         d.Street,
		Alley:          d.Alley,
		ApartmentPlate: d.ApartmentPlate,
		ApartmentNo:    d.ApartmentNo,
		PostalCode:     d.PostalCode,
		LicenseNo:      d.LicenseNo,
		LicenseStatus:  d.LicenseStatus,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}

func NewAuthService(dao repository.DAO, session *scs.SessionManager) AuthService {
	return &authService{dao: dao, session: session}
}

func (a *authService) SignupUser(ctx context.Context, payload dto.CreateUser) (dto.User, error) {
	user, err := a.dao.NewUserRepository().GetUserByPhone(payload.Phone)

	// we shouldn't consider ErrNoRows as an error in this case, because we actually want to insert a new user
	if err != nil && err != sql.ErrNoRows {
		return dto.User{}, err
	}

	if user.ID > 0 {
		if user.Status == models.StatusUserInactive {
			signupOTP := otp.GenerateOTPCode(6)
			fmt.Println("already signup, but otp is: ", signupOTP)

			a.session.Put(ctx, "otpData", dto.PhoneWithOTP{
				Phone: payload.Phone,
				OTP:   signupOTP,
			})

			return dto.User{}, serviceerr.NewServiceErr("need to verify the account", http.StatusBadRequest)
		}

		return dto.User{}, serviceerr.NewServiceErr("user already exists", http.StatusBadRequest)
	}

	hashedPassword, err := password.HashPassword(payload.Password)
	if err != nil {
		return dto.User{}, err
	}

	// insert user, if not already inserted(with default inactive status)
	_, err = a.dao.NewUserRepository().CreateUser(dto.CreateUser{
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
		return dto.User{}, err
	}

	signupOTP := otp.GenerateOTPCode(6)
	fmt.Println("signup otp: ", signupOTP)
	a.session.Put(ctx, "otpData", dto.PhoneWithOTP{
		Phone: payload.Phone,
		OTP:   signupOTP,
	})

	/* TODO: send the validation email(if provided) and SMS, so that user can verify both, but the SMS verification is necessary for the user to be registered.
	We can use the message field of Resp type and make it to have fa and en.*/

	return UserModelToDTO(user), nil
}

func (a *authService) SignupDriver(ctx context.Context, payload dto.SignupDriverRequest) (dto.Driver, error) {
	driver, err := a.dao.NewDriverRepository().GetDriverByPhone(payload.Phone)
	if err != nil && err != sql.ErrNoRows {
		return dto.Driver{}, serviceerr.NewServiceErr("sth went wrong", http.StatusBadRequest)
	}

	if driver.ID > 0 {
		/* Generate an OTP again so user can continue the signup process*/
		if driver.Status == models.StatusDriverInactive {
			signupOTP := otp.GenerateOTPCode(6)
			fmt.Println("already signup, but otp is: ", signupOTP)

			a.session.Put(ctx, "otpData", dto.PhoneWithOTP{
				Phone: payload.Phone,
				OTP:   signupOTP,
			})

			return dto.Driver{}, serviceerr.NewServiceErr("need to verify the account", http.StatusBadRequest)
		}

		return dto.Driver{}, serviceerr.NewServiceErr("driver already exists", http.StatusBadRequest)
	}

	hashedPassword, err := password.HashPassword(payload.Password)
	if err != nil {
		return dto.Driver{}, serviceerr.NewServiceErr("sth went wrong", http.StatusBadRequest)
	}

	driver, err = a.dao.NewDriverRepository().CreateDriver(dto.CreateDriver{
		Phone:          payload.Phone,
		Password:       hashedPassword,
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
	})
	if err != nil {
		return dto.Driver{}, serviceerr.NewServiceErr("sth went wrong", http.StatusBadRequest)
	}

	// TODO: validate license_no using an external API(R&D this) ...

	signupOTP := otp.GenerateOTPCode(6)
	fmt.Println("signup otp: ", signupOTP)
	a.session.Put(ctx, "otpData", dto.PhoneWithOTP{
		Phone: payload.Phone,
		OTP:   signupOTP,
	})

	return driverModelToDTO(driver), nil
}

func (a *authService) VerifyUserSignup(ctx context.Context, payload dto.VerifyUserSignupRequest) error {
	OTPData, ok := a.session.Get(ctx, "otpData").(dto.PhoneWithOTP)
	if !ok {
		return serviceerr.NewServiceErr("please try again", http.StatusBadRequest)
	}

	// check if OTP is correct
	if payload.OTP != OTPData.OTP {
		return serviceerr.NewServiceErr("invalid OTP", http.StatusBadRequest)
	}

	user, err := a.dao.NewUserRepository().GetUserByPhone(OTPData.Phone)
	if err != nil {
		return err
	}

	// update user status to active
	err = a.dao.NewUserRepository().UpdateUserStatus("active", user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *authService) Login(payload dto.LoginRequest) (string, string, error) {
	var u model.User
	var d model.Driver
	var err error

	if payload.IsDriver {
		d, err = a.dao.NewDriverRepository().GetDriverByPhone(payload.Phone)
		if err != nil {
			return "", "", err
		}
	} else {
		u, err = a.dao.NewUserRepository().GetUserByPhone(payload.Phone)
		if err != nil {
			return "", "", serviceerr.NewServiceErr("invalid credentials", http.StatusBadRequest)
		}
	}

	var isPasswordCorrect bool
	if payload.IsDriver {
		isPasswordCorrect = password.CheckPasswordHash(payload.Password, d.Password)
	} else {
		isPasswordCorrect = password.CheckPasswordHash(payload.Password, u.Password)
	}

	if !isPasswordCorrect {
		return "", "", serviceerr.NewServiceErr("invalid credentials", http.StatusBadRequest)
	}

	var accessToken, refreshToken string
	if payload.IsDriver {
		// TODO: Generate JWT and send it back (http cookie or in response?)
		accessToken, refreshToken, err = appjwt.GenerateWithMoreClaims(u.ID, map[string]any{"isDriver": true})
		if err != nil {
			return "", "", err
		}
	} else {
		accessToken, refreshToken, err = appjwt.Generate(u.ID)
		if err != nil {
			return "", "", err
		}
	}

	return accessToken, refreshToken, nil
}
