package model

import "time"

type Pickup struct {
	ID        int
	TruckID   int
	UserID    int
	Time      time.Time
	Weight    float32
	Note      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
