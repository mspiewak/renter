package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/mspiewak/renter/model"
)

type BillRepository struct {
	db *sqlx.DB
}

func NewBillRepository(db *sqlx.DB) *BillRepository {
	return &BillRepository{
		db: db,
	}
}

func (r *BillRepository) GetAll() ([]model.Bill, error) {
	var bills []model.Bill
	err := r.db.Select(&bills, "SELECT id, total_price, due_date, period_start, period_end, url, bill_type_id FROM bill")
	return bills, err
}

func (r *BillRepository) GetByTenantId(id int) ([]model.TenantBill, error) {
	var bills []model.TenantBill
	err := r.db.Select(&bills, `
		SELECT 
			tb.id, tb.price, tb.payment_date, tb.tenant_id, tb.bill_id "bill.id",
			b.total_price "bill.total_price", b.due_date "bill.due_date", b.period_start "bill.period_start", b.period_end "bill.period_end", b.url "bill.url", b.bill_type_id "bill.bill_type.id",
			bt.name "bill.bill_type.name"
		FROM tenant_bill tb 
		LEFT JOIN bill b ON (b.id=tb.bill_id) 
		LEFT JOIN bill_type bt ON (b.bill_type_id=bt.id) 
		WHERE tenant_id = ?
		ORDER BY b.due_date DESC
		`, id)
	return bills, err
}

func (r *BillRepository) CreateTenantBill(tb *model.TenantBill) error {
	res, err := r.db.Exec("INSERT INTO tenant_bill (tenant_id, price, bill_id) VALUES (?,?,?)", tb.TenantID, tb.Price, tb.Bill.ID)
	if err != nil {
		return fmt.Errorf("cannot insert row: %v", err)
	}
	lastInsertID, _ := res.LastInsertId()
	tb.ID = int(lastInsertID)
	return nil
}

func (r *BillRepository) CreateBill(b *model.Bill) error {
	res, err := r.db.Exec(
		"INSERT INTO bill (total_price, due_date, period_start, period_end, url, bill_type_id) VALUES (?,?,?,?,?,?)",
		b.Price, b.DueDate, b.PeriodStart, b.PeriodEnd, b.URL, b.Type.ID,
	)
	if err != nil {
		return fmt.Errorf("cannot insert row: %v", err)
	}
	lastInsertID, _ := res.LastInsertId()
	b.ID = int(lastInsertID)

	return nil
}
