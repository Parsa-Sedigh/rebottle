package repository

import (
	"context"
	"database/sql"
	"github.com/Parsa-Sedigh/rebottle/internal/dto"
	"github.com/Parsa-Sedigh/rebottle/internal/model"
)

type UserRepository interface {
	CreateUser(user dto.CreateUser) (int, error)
	GetUserByPhone(phone string) (model.User, error)
	GetUserByID(id int) (model.User, error)
	UpdateUser(user dto.UpdateUser) (model.User, error)
	UpdateUserStatus(status string, userID int) error
	GetUsers() ([]model.User, error)
	InactiveUser() error
}

type userRepository struct {
	DB *sql.DB
}

func scanUserRow(userRow *sql.Row, u *model.User) error {
	err := userRow.Scan(
		&u.ID,
		&u.Phone,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.Credit,
		&u.Status,
		&u.Province,
		&u.City,
		&u.Street,
		&u.Alley,
		&u.ApartmentPlate,
		&u.ApartmentNo,
		&u.PostalCode,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) CreateUser(user dto.CreateUser) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var id int

	stmt := `INSERT INTO "user"(phone, first_name, last_name, email, password, credit, province, city,
		       street, alley, apartment_plate, apartment_no, postal_code)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id`

	err := u.DB.QueryRowContext(ctx, stmt,
		user.Phone,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		0,
		user.Province,
		user.City,
		user.Street,
		user.Alley,
		user.ApartmentPlate,
		user.ApartmentNo,
		user.PostalCode,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *userRepository) GetUsers() ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var users []model.User

	query := `SELECT id, phone, first_name, last_name, email, password, credit, status, province, city,
			street, alley, apartment_plate, apartment_no, postal_code, created_at, updated_at
			FROM "user"`
	rows, err := u.DB.QueryContext(ctx, query)
	if err != nil {
		return []model.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User

		err = rows.Scan(
			&user.ID,
			&user.Phone,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.Credit,
			&user.Status,
			&user.Province,
			&user.City,
			&user.Street,
			&user.Alley,
			&user.ApartmentPlate,
			&user.ApartmentNo,
			&user.PostalCode,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return []model.User{}, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return []model.User{}, err
	}

	return users, nil
}

func (u *userRepository) GetUserByID(id int) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user model.User

	row := u.DB.QueryRowContext(ctx, `SELECT id, phone, first_name, last_name, email, password, credit, status, province, city,
			street, alley, apartment_plate, apartment_no, postal_code, created_at, updated_at
			FROM "user" WHERE id = $1`, id)
	err := scanUserRow(row, &user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *userRepository) GetUserByPhone(phone string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user model.User

	row := u.DB.QueryRowContext(ctx, `SELECT id, phone, first_name, last_name, email, password, credit, status, province, city,
			street, alley, apartment_plate, apartment_no, postal_code, created_at, updated_at
			FROM "user" WHERE phone = $1`, phone)

	err := scanUserRow(row, &user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *userRepository) UpdateUser(user dto.UpdateUser) (model.User, error) {
	return model.User{}, nil
}

func (u *userRepository) InactiveUser() error {
	//ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	//defer cancel()

	return nil
}

func (u *userRepository) UpdateUserStatus(status string, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		UPDATE "user" SET status = $1 WHERE id = $2
	`
	_, err := u.DB.ExecContext(ctx, stmt, status, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) UpdateUserPassword(hash string, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		UPDATE "user" SET password = $1 WHERE id = $2
	`
	_, err := u.DB.ExecContext(ctx, stmt, hash, userID)
	if err != nil {
		return err
	}

	return nil
}
