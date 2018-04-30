package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mspiewak/renter/model"
)

func (a *App) postBill(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var bill model.Bill
	if err := json.NewDecoder(r.Body).Decode(&bill); err != nil {
		return nil, fmt.Errorf("cannot decode json: %v", err)
	}

	noOfDaysRenting, err := a.tenantRepository.GetNumberOfRentingDays(bill.PeriodStart, bill.PeriodEnd)
	if err != nil {
		return nil, fmt.Errorf("cannot get number of renting days: %v", err)
	}

	if err := a.billRepository.CreateBill(&bill); err != nil {
		return nil, err
	}

	sum := 0
	for _, noOfDays := range noOfDaysRenting {
		sum += noOfDays
	}

	for tenantID, noOfDays := range noOfDaysRenting {
		tb := model.TenantBill{
			TenantID: tenantID,
			Price:    float64(noOfDays) * bill.Price / float64(sum),
			Bill: model.Bill{
				ID: bill.ID,
			},
		}
		if err := a.billRepository.CreateTenantBill(&tb); err != nil {
			return nil, fmt.Errorf("cannot create bill for tenant: %v", err)
		}
	}

	return bill, nil
}

func (a *App) getBills(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return a.billRepository.GetAll()
}

func (a *App) getTenantBills(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	tenantID := getRealTenantID(vars["id"])

	return a.billRepository.GetByTenantId(tenantID)
}
