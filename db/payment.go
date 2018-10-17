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
