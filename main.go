package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type App struct {
	DB *sqlx.DB
}

func main() {
	var app App
	app.Initialize()
	defer app.DB.Close()

	r := mux.NewRouter()
	r.HandleFunc("/tenant", app.getTenantsHandler).Methods(http.MethodGet)
	r.HandleFunc("/tenant/{id:[0-9]+}/bill", app.getTenantBills).Methods(http.MethodGet)
	r.HandleFunc("/bill", app.getBills).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8090", r))
}

func (a *App) Initialize() {
	var err error
	a.DB, err = sqlx.Open("mysql", "root:root@tcp(127.0.0.1:3306)/renter?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	if err = a.DB.Ping(); err != nil {
		log.Fatal(err)
	}
}

func (a *App) getTenantsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tenants := getTenants(a.DB)
	response, err := json.Marshal(tenants)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (a *App) getBills(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bills := getBills(a.DB)
	response, err := json.Marshal(bills)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (a *App) getTenantBills(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	vars := mux.Vars(r)
	tenantID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bills := getTenantBills(a.DB, tenantID)
	response, err := json.Marshal(bills)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
