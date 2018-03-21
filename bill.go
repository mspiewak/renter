package main

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// BillType keeps information about bill type, e.x. electricity, gas
type BillType struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type NestedBillType struct {
	ID int `json:"id" db:"bill_type_id"`
	BillType
}

// Bill keeps information about particular bill
type Bill struct {
	ID             int       `json:"id" db:"id"`
	Price          float32   `json:"total_price" db:"total_price"`
	DueDate        time.Time `json:"due_date" db:"due_date"`
	PeriodStart    time.Time `json:"period_start" db:"period_start"`
	PeriodEnd      time.Time `json:"period_end" db:"period_end"`
	URL            *string   `json:"url" db:"url"`
	NestedBillType `json:"type"`
}

type NestedBill struct {
	ID int `json:"id" db:"bill_id"`
	Bill
}

// TenantBill keeps information about bill breakdown per particular tenant
type TenantBill struct {
	ID          int        `json:"id" db:"tid"`
	TenantID    int        `json:"tenant_id" db:"tenant_id"`
	Price       float32    `json:"price" db:"price"`
	PaymentDate *time.Time `json:"payment_date" db:"payment_date"`
	NestedBill  `json:"bill"`
}

func getBills(db *sqlx.DB) ([]Bill, error) {
	var bills []Bill
	err := db.Select(&bills, "SELECT * FROM bill")
	return bills, err
}

func getTenantBills(db *sqlx.DB, id int) ([]TenantBill, error) {

	var bills []TenantBill
	err := db.Select(&bills, `
		SELECT 
			tb.id as tid, tb.price, tb.payment_date, tb.tenant_id, tb.bill_id,
			b.total_price, b.due_date, b.period_start, b.period_end, b.url, b.bill_type_id, b.url,
			bt.name
		FROM tenant_bill tb 
		LEFT JOIN bill b ON (b.id=tb.bill_id) 
		LEFT JOIN bill_type bt ON (b.bill_type_id=bt.id) 
		WHERE tenant_id = ?
		ORDER BY b.due_date DESC
		`, id)
	return bills, err
}
