package model

import "time"

// Payment keeps information about particular payment transaction
type Payment struct {
	ID          int        `json:"id" db:"id"`
	Amount      float64    `json:"amount" db:"amount"`
	PaymentDate *time.Time `json:"payment_date" db:"payment_date"`
	TenantID    int        `json:"tenant_id" db:"tenant_id"`
}
