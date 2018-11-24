package main

import (
	"net/http"
)

const taxRate = 8.5

func (a *App) getIncome(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return a.paymentRepository.GetIncome(taxRate)
}
