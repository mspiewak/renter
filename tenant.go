package main

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

// Tenant keeps information about particular tenants
type Tenant struct {
	ID          int        `json:"id" db:"id"`
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	MoveInDate  *time.Time `json:"move_in_date" db:"move_in_date"`
	MoveOutDate *time.Time `json:"move_out_date" db:"move_out_date"`
	RoomID      int        `json:"room_id" db:"room_id"`
	Password    string     `json:"-" db:"password"`
}

func getTenants(db *sqlx.DB) []Tenant {
	var tenants []Tenant
	if err := db.Select(&tenants, "SELECT * FROM tenant"); err != nil {
		log.Fatal(err)
	}
	return tenants
}
