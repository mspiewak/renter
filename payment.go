package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) getTenantPayments(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	tenantID := getRealTenantID(vars["id"])

	return a.paymentRepository.GetByTenantId(tenantID)
}
