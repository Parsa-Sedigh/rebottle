package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Parsa-Sedigh/rebottle/internal/password"
	"time"
)

type Models struct {
	DB *sql.DB
}

func NewModels(db *sql.DB) Models {
	return Models{DB: db}
}

const (
	StatusPickupWaiting           = "waiting"
	StatusPickupInProcess         = "in_process"
	StatusPickupCancelledByUser   = "cancelled_by_user"
	StatusPickupCancelledBySystem = "cancelled_by_system"
	StatusPickupDone              = "done"

	StatusUserEmailInactive = "inactive"
	StatusUserEmailActive   = "active"
)

type User struct {
	ID             int       `json:"id"`
	Phone          string    `json:"phone"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Password       string    `json:"-"`
	Credit         uint16    `json:"credit"`
	Status         string    `json:"status"` // TODO: how convert sql enums to go code?
	EmailStatus    string    `json:"email_status"`
	Province       string    `json:"province"`
	City           string    `json:"city"`
	Street         string    `json:"street"`
	Alley          string    `json:"alley"`
	ApartmentPlate uint16    `json:"apartment_plate"`
	ApartmentNo    uint16    `json:"apartment_no"`
	PostalCode     string    `json:"postal_code"`
	CreatedAt      time.Time `json:"created_at"` // TODO: how convert sql timestamp to go code?
	UpdatedAt      time.Time `json:"updated_at"`
}

type SignupUserRequest struct {
	Phone          string `json:"phone" validate:"required,min=11,max=11,phone"`
	FirstName      string `json:"first_name" validate:"required,min=3"`
	LastName       string `json:"last_name" validate:"required,min=3"`
	Email          string `json:"email,omitempty" validate:"omitempty,email"`
	Password       string `json:"password" validate:"min=6"`
	Province       string `json:"province" validate:"required"`
	City           string `json:"city" validate:"required"`
	Street         string `json:"street" validate:"required"`
	Alley          string `json:"alley,omitempty"`
	ApartmentPlate int    `json:"apartment_plate,omitempty" validate:"required"`
	ApartmentNo    uint16 `json:"apartment_no,omitempty" validate:"required"`
	PostalCode     string `json:"postal_code" validate:"required"`
}

type Pickup struct {
	ID        int       `json:"id"`
	TruckID   int       `json:"truck_id"`
	UserID    int       `json:"user_id"`
	Time      time.Time `json:"time"`
	Weight    float32   `json:"weight"`
	Note      string    `json:"note"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Driver struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	LicenseNo string    `json:"license_no"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Truck struct {
	ID        int       `json:"id"`
	PlateNo   string    `json:"plate_no"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func genGetUserByFieldQuery(field string) string {
	return fmt.Sprintf(`SELECT id, phone, first_name, last_name, email, password, credit, status, province, city,
			street, alley, apartment_plate, apartment_no, postal_code, created_at, updated_at
			FROM "user" WHERE %s = $1`, field)
}

func (m *Models) scanUserRow(userRow *sql.Row, u *User) error {
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

func (m *Models) GetUserByPhone(phone string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u User

	row := m.DB.QueryRowContext(ctx, genGetUserByFieldQuery("phone"), phone)

	err := m.scanUserRow(row, &u)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (m *Models) GetUserByID(id int) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u User

	row := m.DB.QueryRowContext(ctx, genGetUserByFieldQuery("id"), id)
	err := m.scanUserRow(row, &u)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (m *Models) GetUserByEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u User

	row := m.DB.QueryRowContext(ctx, genGetUserByFieldQuery("id"), email)
	err := m.scanUserRow(row, &u)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (m *Models) InsertUser(u SignupUserRequest) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	hashedPassword, err := password.HashPassword(u.Password)
	if err != nil {
		return 0, err
	}

	var id int

	stmt := `
		INSERT INTO "user"(phone, first_name, last_name, email, password, credit, province, city,
		       street, alley, apartment_plate, apartment_no, postal_code)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`
	_, err = m.DB.ExecContext(ctx, stmt, u.Phone, u.FirstName, u.LastName, u.Email, hashedPassword, 0, u.Province, u.City,
		u.Street, u.Alley, u.ApartmentPlate, u.ApartmentNo, u.PostalCode)
	if err != nil {
		return 0, err
	}

	stmt = `SELECT id FROM "user" WHERE phone = $1`
	err = m.DB.QueryRowContext(ctx, stmt, u.Phone).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *Models) UpdateUserStatus(status string, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		UPDATE "user" SET status = $1 WHERE id = $2
	`
	_, err := m.DB.ExecContext(ctx, stmt, status, userID)
	if err != nil {
		return err
	}

	return nil
}

func (m *Models) GetPickup(pickupID, userID int) (Pickup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var p Pickup

	/* Users can only see their own pickups, so we should check the userID too. */
	err := m.DB.QueryRowContext(ctx, `
		SELECT id, COALESCE(truck_id, 0), user_id, weight, "time", note, status, created_at, updated_at
		FROM pickup 
		WHERE id = $1 AND user_id = $2`, pickupID, userID).Scan(
		&p.ID,
		&p.TruckID,
		&p.UserID,
		&p.Weight,
		&p.Time,
		&p.Note,
		&p.Status,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return Pickup{}, err
	}

	return p, nil
}

func (m *Models) InsertPickup(p Pickup) (Pickup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO pickup(truck_id, user_id, weight, time, note, status)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, COALESCE(truck_id, 0), user_id, weight, time, note, status, created_at, updated_at
	`

	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return Pickup{}, err
	}
	defer stmt.Close()

	var insertedPickup Pickup

	// TODO: how can I use exec here but also get the lastInsertID() (knowing that we can't use lastInsertID() in pgx). Also I can't do a query after this either
	err = stmt.QueryRowContext(ctx, sql.NullInt16{Int16: int16(p.TruckID)}, p.UserID, p.Weight, p.Time, p.Note, p.Status).Scan(
		&insertedPickup.ID,
		&insertedPickup.TruckID,
		&insertedPickup.UserID,
		&insertedPickup.Weight,
		&insertedPickup.Time,
		&insertedPickup.Note,
		&insertedPickup.Status,
		&insertedPickup.CreatedAt,
		&insertedPickup.UpdatedAt,
	)
	if err != nil {
		return Pickup{}, err
	}

	return insertedPickup, nil
}

type UpdatePickupParams struct {
	ID     int
	UserID int
	Time   time.Time
	Weight float32
	Note   string
	Status string
}

func (m *Models) UpdatePickup(params UpdatePickupParams) (Pickup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		UPDATE pickup SET weight = $1, time = $2, note = $3, status = $4 WHERE id = $5
	`
	_, err := m.DB.ExecContext(ctx, stmt, params.Weight, params.Time, params.Note, params.Status, params.ID)
	if err != nil {
		return Pickup{}, err
	}

	p, err := m.GetPickup(params.ID, params.UserID)
	if err != nil {
		return Pickup{}, err
	}

	return p, nil
}

func (m *Models) GetUserPickups(userID int) ([]Pickup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var pickups []Pickup

	rows, err := m.DB.QueryContext(ctx, `
		SELECT id, COALESCE(truck_id, 0), user_id, weight, time, note, status, created_at, updated_at
		FROM pickup
		WHERE user_id = $1`, userID)
	if err != nil {
		return pickups, nil
	}
	defer rows.Close()

	for rows.Next() {
		var p Pickup
		if err = rows.Scan(
			&p.ID,
			&p.TruckID,
			&p.UserID,
			&p.Weight,
			&p.Time,
			&p.Note,
			&p.Status,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return pickups, err
		}

		pickups = append(pickups, p)
	}

	if err = rows.Err(); err != nil {
		return pickups, err
	}
	if err = rows.Close(); err != nil {
		return pickups, err
	}

	return pickups, nil
}

func (m *Models) CancelPickup(pickupID, userID int, byUser bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var status string
	if byUser {
		status = StatusPickupCancelledByUser
	} else {
		status = StatusPickupCancelledBySystem
	}

	stmt := `UPDATE pickup SET status = $1 WHERE id = $2 AND user_id = $3`
	_, err := m.DB.ExecContext(ctx, stmt, status, pickupID, userID)
	if err != nil {
		return err
	}

	return nil
}
