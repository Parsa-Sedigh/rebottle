package models

import (
	"context"
	"database/sql"
	"fmt"
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

	StatusUserInactive = "inactive"
	StatusUserBlocked  = "blocked"
	StatusUserActive   = "active"

	StatusDriverInactive = "inactive"
	StatusDriverBlocked  = "blocked"
	StatusDriverActive   = "active"

	StatusUserEmailInactive = "inactive"
	StatusUserEmailActive   = "active"
)

const dbTimeout = 3 * time.Second

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
	ApartmentNo    int    `json:"apartment_no,omitempty" validate:"required"`
	PostalCode     string `json:"postal_code" validate:"required"` // TODO: Better validation
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
	ID             int    `json:"id"`
	Phone          string `json:"phone"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Password       string `json:"-"`
	Status         string `json:"status"` // TODO: how convert sql enums to go code?
	EmailStatus    string `json:"email_status"`
	Province       string `json:"province"`
	City           string `json:"city"`
	Street         string `json:"street"`
	Alley          string `json:"alley"`
	ApartmentPlate uint16 `json:"apartment_plate"`
	ApartmentNo    uint16 `json:"apartment_no"`
	PostalCode     string `json:"postal_code"`
	//UserID    int       `json:"user_id"`
	LicenseNo     string    `json:"license_no"`
	LicenseStatus string    `json:"license_status"`
	CreatedAt     time.Time `json:"created_at"` // TODO: how convert sql timestamp to go code?
	UpdatedAt     time.Time `json:"updated_at"`
}

type InsertDriverData struct {
	Phone          string `json:"phone"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	LicenseNo      string `json:"license_no"`
	Province       string `json:"province"`
	City           string `json:"city"`
	Street         string `json:"street"`
	Alley          string `json:"alley"`
	ApartmentPlate uint16 `json:"apartment_plate"`
	ApartmentNo    uint16 `json:"apartment_no"`
	PostalCode     string `json:"postal_code"`
}

type Truck struct {
	ID        int       `json:"id"`
	PlateNo   string    `json:"plate_no"`
	Model     string    `json:"model"`
	Color     string    `json:"color"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func genGetUserByFieldQuery(field string) string {
	return fmt.Sprintf(`SELECT id, phone, first_name, last_name, email, password, credit, status, province, city,
			street, alley, apartment_plate, apartment_no, postal_code, created_at, updated_at
			FROM "user" WHERE %s = $1`, field)
}

func genGetDriverByFieldQuery(field string) string {
	return fmt.Sprintf(`SELECT id, phone, first_name, last_name, email, license_no, license_status, status, email_status, province, city, street, alley,
	apartment_plate, apartment_no, postal_code, created_at, updated_at FROM driver WHERE %s = $1`, field)
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

func (m *Models) scanDriverRow(driverRow *sql.Row, d *Driver) error {
	err := driverRow.Scan(
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
	fmt.Println("scan err: ", err)
	if err != nil {
		return err
	}

	return nil
}

func (m *Models) GetUserByPhone(phone string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var id int

	stmt := `
		INSERT INTO "user"(phone, first_name, last_name, email, password, credit, province, city,
		       street, alley, apartment_plate, apartment_no, postal_code)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`
	_, err := m.DB.ExecContext(ctx, stmt, u.Phone, u.FirstName, u.LastName, u.Email, u.Password, 0, u.Province, u.City,
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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
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

func (m *Models) UpdateUserPassword(hash string, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		UPDATE "user" SET password = $1 WHERE id = $2
	`
	_, err := m.DB.ExecContext(ctx, stmt, hash, userID)
	if err != nil {
		return err
	}

	return nil
}

func (m *Models) GetPickup(pickupID, userID int) (Pickup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
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

func (m *Models) InsertDriver(data InsertDriverData) (Driver, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		INSERT INTO driver(phone, first_name, last_name, email, password, license_no, province, city,
			street, alley, apartment_plate, apartment_no, postal_code)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, phone, first_name, last_name, email, license_no, license_status, status, email_status,
		province, city, street, alley, apartment_plate, apartment_no, postal_code, created_at, updated_at
	`
	preparedStmt, err := m.DB.PrepareContext(ctx, stmt)
	if err != nil {
		return Driver{}, err
	}
	defer preparedStmt.Close()

	var driver Driver

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
	err = m.scanDriverRow(row, &driver)
	if err != nil {
		return Driver{}, err
	}

	return driver, nil
}

func (m *Models) GetDriverByID(id int) (Driver, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var d Driver

	row := m.DB.QueryRowContext(ctx, genGetDriverByFieldQuery("id"), id)
	err := m.scanDriverRow(row, &d)
	if err != nil {
		return Driver{}, err
	}

	return d, nil
}

func (m *Models) GetDriverByPhone(phone string) (Driver, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var d Driver

	row := m.DB.QueryRowContext(ctx, genGetDriverByFieldQuery("phone"), phone)
	err := m.scanDriverRow(row, &d)
	if err != nil {
		return Driver{}, err
	}

	return d, nil
}

type UpdateDriverData struct {
	ID             int
	FirstName      string
	LastName       string
	Email          string
	LicenseNo      string
	Province       string
	City           string
	Street         string
	Alley          string
	ApartmentPlate uint16
	ApartmentNo    uint16
	PostalCode     string
}

func (m *Models) UpdateDriver(data UpdateDriverData) (Driver, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var d Driver

	stmt := `UPDATE driver SET first_name = $1, last_name = $2, email = $3, license_no = $4,
                  province = $5, city = $6, street = $7, alley = $8, apartment_plate = $9, apartment_no = $10, postal_code = $11
                  WHERE id = $12`
	row := m.DB.QueryRowContext(ctx, stmt,
		data.FirstName,
		data.LastName,
		data.Email,
		data.LicenseNo,
		data.Province,
		data.City,
		data.Street,
		data.Alley,
		data.ApartmentPlate,
		data.ApartmentNo,
		data.PostalCode,
		data.ID,
	)
	err := m.scanDriverRow(row, &d)
	if err != nil {
		return Driver{}, err
	}

	return d, nil
}

func (m *Models) UpdateDriverStatus(status string, driverID int) (Driver, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `UPDATE driver SET status = $1 WHERE id = $2`
	preparedStmt, err := m.DB.PrepareContext(ctx, stmt)
	if err != nil {
		return Driver{}, err
	}

	preparedStmt.QueryRowContext(ctx, status, driverID)
	if err != nil {
		return Driver{}, err
	}

	driver, err := m.GetDriverByID(driverID)
	if err != nil {
		return Driver{}, err
	}

	return driver, nil
}
