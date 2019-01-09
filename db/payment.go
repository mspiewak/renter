package db

import (
	"github.com/jmoiron/sqlx"

	"github.com/mspiewak/renter/model"
)

type PaymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) GetByTenantId(id int) ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.Select(&payments, `
		SELECT id, amount, payment_date, tenant_id
		FROM payment 
		WHERE tenant_id = ?
		ORDER BY payment_date DESC
		`, id)
	return payments, err
}

func (r *PaymentRepository) GetIncome(taxRate float32) ([]model.Income, error) {
	var income []model.Income
	err := r.db.Select(&income, `
		SELECT YEAR(due_date) as year, MONTH(due_date) as month, ROUND(SUM(total_price), 2) as total, ROUND(SUM(total_price) * 8.5 / 100, 2) as tax
		FROM bill
		WHERE bill_type_id=1
		GROUP BY YEAR(due_date), MONTH(due_date)
		ORDER BY YEAR(due_date) DESC, MONTH(due_date) DESC
		`, taxRate)
	return income, err
}
