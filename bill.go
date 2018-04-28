package main

import (
	"fmt"
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

func postBill(db *sqlx.DB, b *Bill) error {
	res, err := db.Exec(
		"INSERT INTO bill (total_price, due_date, period_start, period_end, url, bill_type_id) VALUES (?,?,?,?,?,?)",
		b.Price, b.DueDate, b.PeriodStart, b.PeriodEnd, b.URL, b.NestedBillType.ID,
	)
	if err != nil {
		return fmt.Errorf("cannot insert row: %v", err)
	}
	lastInsertID, _ := res.LastInsertId()
	b.ID = int(lastInsertID)

	dd, err := getDaysDistribution(db, *b)
	if err != nil {
		return fmt.Errorf("cannot get days distribution: %v", err)
	}

	sum := 0
	for _, r := range dd {
		sum += r.Days
	}

	for _, r := range dd {
		tb := TenantBill{
			TenantID: r.TenantID,
			Price:    float32(r.Days) * b.Price / float32(sum),
			NestedBill: NestedBill{
				ID: b.ID,
			},
		}
		if err := postTenantBill(db, &tb); err != nil {
			return fmt.Errorf("cannot create bill for tenant: %v", err)
		}
	}

	return nil
}

func postTenantBill(db *sqlx.DB, tb *TenantBill) error {
	res, err := db.Exec("INSERT INTO tenant_bill (tenant_id, price, bill_id) VALUES (?,?,?)", tb.TenantID, tb.Price, tb.NestedBill.ID)
	if err != nil {
		return fmt.Errorf("cannot insert row: %v", err)
	}
	lastInsertID, _ := res.LastInsertId()
	tb.ID = int(lastInsertID)
	return nil
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

type BillDaysDistribution []struct {
	TenantID int `db:"tenant_id"`
	Days     int `db:"days"`
}

func getDaysDistribution(db *sqlx.DB, b Bill) (BillDaysDistribution, error) {
	var d BillDaysDistribution
	nstmt, err := db.PrepareNamed(`
		SELECT 
			tenant.id as tenant_id, 
			DATEDIFF(
				IF(move_out_date < :period_end, move_out_date, :period_end), 
				IF(move_in_date > :period_start, move_in_date, :period_start)
			) as days
		FROM tenant 
		WHERE tenant.move_in_date < :period_end
		AND (tenant.move_out_date IS NULL || tenant.move_out_date > :period_start)
		`)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare statement: %v", err)
	}
	if err := nstmt.Select(&d, b); err != nil {
		return nil, fmt.Errorf("cannot get data: %v", err)
	}
	return d, nil
}
