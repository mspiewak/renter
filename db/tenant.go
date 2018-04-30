package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/mspiewak/renter/model"
)

type TenantRepository struct {
	db *sqlx.DB
}

func NewTenantRepository(db *sqlx.DB) *TenantRepository {
	return &TenantRepository{
		db: db,
	}
}

func (r *TenantRepository) GetAll() ([]model.Tenant, error) {
	var tenants []model.Tenant
	err := r.db.Select(
		&tenants,
		"SELECT id, first_name, last_name, move_in_date, room_id, deposit, rent, move_out_date FROM tenant")
	return tenants, err
}

func (r *TenantRepository) GetByID(id int) (model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.Get(
		&tenant,
		`
		SELECT id, first_name, last_name, move_in_date, room_id, deposit, rent, move_out_date 
		FROM tenant
		WHERE id = ?
		`, id)
	return tenant, err
}

func (r *TenantRepository) GetNumberOfRentingDays(periodStart, periodEnd time.Time) (map[int]int, error) {
	nstmt, err := r.db.PrepareNamed(`
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

	rows, err := nstmt.Queryx(map[string]interface{}{
		"period_start": periodStart,
		"period_end":   periodEnd,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot execute query: %v", err)
	}

	res := make(map[int]int)
	for rows.Next() {
		var tenantID int
		var noOfDays int
		if err := rows.Scan(&tenantID, &noOfDays); err != nil {
			return nil, fmt.Errorf("cannot scan into map: %v", err)
		}
		res[tenantID] = noOfDays
	}

	return res, nil
}
