package repository

import (
	"context"
	"database/sql"
	"github.com/Parsa-Sedigh/rebottle/internal/dto"
	"github.com/Parsa-Sedigh/rebottle/internal/model"
)

type PickupRepository interface {
	CreatePickup(pickup dto.Pickup) (model.Pickup, error)
	GetUserPickups(userID int) ([]model.Pickup, error)
	GetUserPickup(id int, userID int) (model.Pickup, error)
}

type pickupRepository struct {
	DB *sql.DB
}

func (p *pickupRepository) CreatePickup(pickup dto.Pickup) (model.Pickup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO pickup(truck_id, user_id, weight, time, note, status)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, COALESCE(truck_id, 0), user_id, weight, time, note, status, created_at, updated_at
	`

	stmt, err := p.DB.PrepareContext(ctx, query)
	if err != nil {
		return model.Pickup{}, err
	}
	defer stmt.Close()

	var insertedPickup model.Pickup

	// TODO: how can I use exec here but also get the lastInsertID() (knowing that we can't use lastInsertID() in pgx). Also I can't do a query after this either
	err = stmt.QueryRowContext(ctx, sql.NullInt16{Int16: int16(pickup.TruckID)}, pickup.UserID, pickup.Weight, pickup.Time, pickup.Note, pickup.Status).Scan(
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
		return model.Pickup{}, err
	}

	return insertedPickup, nil
}

func (p *pickupRepository) GetUserPickups(userID int) ([]model.Pickup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var pickups []model.Pickup

	rows, err := p.DB.QueryContext(ctx, `
		SELECT id, COALESCE(truck_id, 0), user_id, weight, time, note, status, created_at, updated_at
		FROM pickup
		WHERE user_id = $1`, userID)
	if err != nil {
		return pickups, nil
	}
	defer rows.Close()

	for rows.Next() {
		var pickup model.Pickup
		if err = rows.Scan(
			&pickup.ID,
			&pickup.TruckID,
			&pickup.UserID,
			&pickup.Weight,
			&pickup.Time,
			&pickup.Note,
			&pickup.Status,
			&pickup.CreatedAt,
			&pickup.UpdatedAt,
		); err != nil {
			return pickups, err
		}

		pickups = append(pickups, pickup)
	}

	if err = rows.Err(); err != nil {
		return pickups, err
	}
	if err = rows.Close(); err != nil {
		return pickups, err
	}

	return pickups, nil
}

func (p *pickupRepository) GetUserPickup(id int, userID int) (model.Pickup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var pickup model.Pickup

	/* Users can only see their own pickups, so we should check the userID too. */
	err := p.DB.QueryRowContext(ctx, `
		SELECT id, COALESCE(truck_id, 0), user_id, weight, "time", note, status, created_at, updated_at
		FROM pickup 
		WHERE id = $1 AND user_id = $2`, id, userID).Scan(
		&pickup.ID,
		&pickup.TruckID,
		&pickup.UserID,
		&pickup.Weight,
		&pickup.Time,
		&pickup.Note,
		&pickup.Status,
		&pickup.CreatedAt,
		&pickup.UpdatedAt,
	)
	if err != nil {
		return model.Pickup{}, err
	}

	return pickup, nil
}

func (p *pickupRepository) InsertPickup(pickup dto.Pickup) (model.Pickup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		INSERT INTO pickup(truck_id, user_id, weight, time, note, status)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, COALESCE(truck_id, 0), user_id, weight, time, note, status, created_at, updated_at
	`

	stmt, err := p.DB.PrepareContext(ctx, query)
	if err != nil {
		return model.Pickup{}, err
	}
	defer stmt.Close()

	var insertedPickup model.Pickup

	// TODO: how can I use exec here but also get the lastInsertID() (knowing that we can't use lastInsertID() in pgx). Also I can't do a query after this either
	err = stmt.QueryRowContext(ctx, sql.NullInt16{Int16: int16(pickup.TruckID)}, pickup.UserID, pickup.Weight, pickup.Time, pickup.Note, pickup.Status).Scan(
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
		return model.Pickup{}, err
	}

	return insertedPickup, nil
}
