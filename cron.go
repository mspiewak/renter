package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mspiewak/renter/model"
)

func (a *App) postRent(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	now := time.Now().Add(time.Hour * 24)
	currentYear, currentMonth, _ := now.Date()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, now.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	tenants, err := a.tenantRepository.GetNumberOfRentingDays(firstOfMonth, lastOfMonth)
	if err != nil {
		return nil, fmt.Errorf("cannot get tenants for current month: %v", err)
	}

	var bills []model.Bill
	for tenantID, days := range tenants {
		days++
		tenant, err := a.tenantRepository.GetByID(tenantID)
		if err != nil {
			return nil, fmt.Errorf("cannot get tenant data: %v", err)
		}

		rentVal := tenant.Rent
		if days < lastOfMonth.Day() {
			rentVal = float64(days) * rentVal / float64(lastOfMonth.Day())
		}
		b := model.Bill{
			DueDate:     time.Date(currentYear, currentMonth, 10, 0, 0, 0, 0, now.Location()),
			PeriodStart: firstOfMonth,
			PeriodEnd:   lastOfMonth,
			Type: model.BillType{
				ID: 1,
			},
			Price: rentVal,
		}

		if err := a.billRepository.CreateBill(&b); err != nil {
			return nil, fmt.Errorf("cannot save bill to db: %v", err)
		}

		tb := model.TenantBill{
			TenantID: tenantID,
			Price:    b.Price,
			Bill:     b,
		}

		if err := a.billRepository.CreateTenantBill(&tb); err != nil {
			return nil, fmt.Errorf("cannot save tenant bill to db: %v", err)
		}

		bills = append(bills, b)
	}

	return bills, nil
}
