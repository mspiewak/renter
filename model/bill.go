package model

import "time"

// BillType keeps information about bill type, e.x. electricity, gas
type BillType struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// Bill keeps information about particular bill
type Bill struct {
	ID          int       `json:"id" db:"id"`
	Price       float32   `json:"total_price" db:"total_price"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
	PeriodStart time.Time `json:"period_start" db:"period_start"`
	PeriodEnd   time.Time `json:"period_end" db:"period_end"`
	URL         *string   `json:"url" db:"url"`
	Type        BillType  `json:"type" db:"bill_type"`
}

// TenantBill keeps information about particular tenant bill
type TenantBill struct {
	ID          int        `json:"id" db:"id"`
	TenantID    int        `json:"tenant_id" db:"tenant_id"`
	Price       float32    `json:"price" db:"price"`
	PaymentDate *time.Time `json:"payment_date" db:"payment_date"`
	Bill        `json:"bill" db:"bill"`
}
