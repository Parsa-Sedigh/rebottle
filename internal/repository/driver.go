package repository

import (
	"context"
	"database/sql"
	"github.com/Parsa-Sedigh/rebottle/internal/dto"
	"github.com/Parsa-Sedigh/rebottle/internal/model"
)

type DriverRepository interface {
	GetDriverByPhone(phone string) (model.Driver, error)
	CreateDriver(data dto.CreateDriver) (model.Driver, error)
}

type driverRepository struct {
	DB *sql.DB
}

func scanDriverRow(driverRow *sql.Row, d *model.Driver) error {
	return driverRow.Scan(
		&d.ID,
		&d.Phone,
		&d.FirstName,
		&d.LastName,
		&d.Email,
		&d.LicenseNo,
		&d.LicenseStatus,
		&d.Status,
		&d.EmailStatus,
		&d.Province,
		&d.City,
		&d.Street,
		&d.Alley,
		&d.ApartmentPlate,
		&d.ApartmentNo,
		&d.PostalCode,
		&d.CreatedAt,
		&d.UpdatedAt,
	)
}

func (d *driverRepository) GetDriverByPhone(phone string) (model.Driver, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var driver model.Driver

	row := d.DB.QueryRowContext(ctx, `SELECT id, phone, first_name, last_name, email, license_no, license_status,
       status, email_status, province, city, street, alley, apartment_plate, apartment_no, postal_code, created_at, updated_at 
		FROM driver WHERE phone = $1`, phone)
	err := scanDriverRow(row, &driver)
	if err != nil {
		return model.Driver{}, err
	}

	return driver, nil
}

func (d *driverRepository) CreateDriver(data dto.CreateDriver) (model.Driver, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		INSERT INTO driver(phone, first_name, last_name, email, password, license_no, province, city,
			street, alley, apartment_plate, apartment_no, postal_code)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, phone, first_name, last_name, email, license_no, license_status, status, email_status,
		province, city, street, alley, apartment_plate, apartment_no, postal_code, created_at, updated_at
	`
	preparedStmt, err := d.DB.PrepareContext(ctx, stmt)
	if err != nil {
		return model.Driver{}, err
	}
	defer preparedStmt.Close()

	var driver model.Driver

	row := preparedStmt.QueryRowContext(ctx,
		data.Phone,
		data.FirstName,
		data.LastName,
		data.Email,
		data.Password,
		data.LicenseNo,
		data.Province,
		data.City,
		data.Street,
		data.Alley,
		data.ApartmentPlate,
		data.ApartmentNo,
		data.PostalCode)
	err = scanDriverRow(row, &driver)
	if err != nil {
		return model.Driver{}, err
	}

	return driver, nil
}
