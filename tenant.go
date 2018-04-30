package main

import "net/http"

func (a *App) getTenantsHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return a.tenantRepository.GetAll()
}
