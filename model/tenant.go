package model

import "time"

// Tenant keeps information about particular tenants
type Tenant struct {
	ID          int        `json:"id" db:"id"`
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	MoveInDate  *time.Time `json:"move_in_date" db:"move_in_date"`
	MoveOutDate *time.Time `json:"move_out_date" db:"move_out_date"`
	RoomID      int        `json:"room_id" db:"room_id"`
	Deposit     float64    `json:"deposit" db:"deposit"`
	Rent        float64    `json:"rent" db:"rent"`
	Password    string     `json:"-" db:"password"`
}
