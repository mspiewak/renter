package main

import (
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
	r.Handle("/tenant", errorHandler(app.getTenantsHandler)).Methods(http.MethodGet)
	r.Handle("/tenant/{id:[0-9]+}/bill", errorHandler(app.getTenantBills)).Methods(http.MethodGet)
	r.Handle("/bill", errorHandler(app.getBills)).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8090", commonHeaders(r)))
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

func (a *App) getTenantsHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return getTenants(a.DB)
}

func (a *App) getBills(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return getBills(a.DB)
}

func (a *App) getTenantBills(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	tenantID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return nil, err
	}

	return getTenantBills(a.DB, tenantID)
}
