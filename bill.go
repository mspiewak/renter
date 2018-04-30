package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mspiewak/renter/model"
)

type BillDaysDistribution []struct {
	TenantID int `db:"tenant_id"`
	Days     int `db:"days"`
}

func getDaysDistribution(db *sqlx.DB, b model.Bill) (BillDaysDistribution, error) {
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
