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
	BillType
	ID int `json:"id" db:"bill_type_id"`
}

// Bill keeps information about particular bill
type Bill struct {
	NestedBillType `json:"type"`
	ID             int       `json:"id" db:"id"`
	Price          float32   `json:"price" db:"price"`
	DueDate        time.Time `json:"due_date" db:"due_date"`
	PeriodStart    time.Time `json:"period_start" db:"period_start"`
	PeriodEnd      time.Time `json:"period_end" db:"period_end"`
	URL            *string   `json:"url" db:"url"`
}

type NestedBill struct {
	Bill
	ID int `json:"id" db:"bill_id"`
}

// TenantBill keeps information about bill breakdown per particular tenant
type TenantBill struct {
	NestedBill  `json:"bill"`
	ID          int        `json:"id" db:"id"`
	TenantID    int        `json:"tenant_id" db:"tenant_id"`
	Price       float32    `json:"price" db:"price"`
	PaymentDate *time.Time `json:"payment_date" db:"payment_date"`
}

func getBills(db *sqlx.DB) ([]Bill, error) {
	var bills []Bill
	err := db.Select(&bills, "SELECT * FROM bill")
	return bills, err
}

func getTenantBills(db *sqlx.DB, id int) ([]TenantBill, error) {
	var bills []TenantBill
	err := db.Select(&bills, `
		SELECT *
		FROM tenant_bill tb 
		LEFT JOIN bill b ON (b.id=tb.bill_id) 
		LEFT JOIN bill_type bt ON (b.bill_type_id=bt.id) 
		WHERE tenant_id = ?
		`, id)
	return bills, err
}
