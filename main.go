package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/mspiewak/renter/db"
	"github.com/mspiewak/renter/model"
)

// App keeps application dependencies
type App struct {
	DB             *sqlx.DB
	billRepository *db.BillRepository
}

func main() {
	var app App

	for {
		log.Println("connecting")
		err := app.Initialize()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 10)
			continue
		}

		break
	}

	defer app.DB.Close()

	log.Println("succesfully connected")

	r := mux.NewRouter()
	r.Handle("/tenant", errorHandler(app.getTenantsHandler)).Methods(http.MethodGet)
	r.Handle("/tenant/{id:[0-9a-z]+}/bill", errorHandler(app.getTenantBills)).Methods(http.MethodGet)
	r.Handle("/bill", errorHandler(app.getBills)).Methods(http.MethodGet)
	r.Handle("/bill", errorHandler(app.postBill)).Methods(http.MethodPost)
	r.Handle("/", errorHandler(app.getBills)).Methods(http.MethodGet)
	r.Handle("/tet", errorHandler(app.getBills)).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8090", commonHeaders(r)))
}

// Initialize the application
func (a *App) Initialize() error {
	dbc, err := sqlx.Connect("mysql", "root:root@tcp(localhost:3306)/renter?parseTime=true")
	if err != nil {
		return fmt.Errorf("cannot connect to db server: %v", err)
	}

	a.DB = dbc
	a.billRepository = db.NewBillRepository(dbc)
	return nil
}

func (a *App) getTenantsHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return getTenants(a.DB)
}

func (a *App) postBill(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var bill model.Bill
	if err := json.NewDecoder(r.Body).Decode(&bill); err != nil {
		return nil, fmt.Errorf("cannot decode json: %v", err)
	}

	if err := a.billRepository.CreateBill(&bill); err != nil {
		return nil, err
	}

	dd, err := getDaysDistribution(a.DB, bill)
	if err != nil {
		return nil, fmt.Errorf("cannot get days distribution: %v", err)
	}

	sum := 0
	for _, r := range dd {
		sum += r.Days
	}

	for _, r := range dd {
		tb := model.TenantBill{
			TenantID: r.TenantID,
			Price:    float32(r.Days) * bill.Price / float32(sum),
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

func getRealTenantID(hash string) int {
	switch hash {
	case "i4lrehq":
		return 1
	case "ggl0qk8":
		return 2
	case "4duspw0":
		return 3
	case "wm4yk48":
		return 4
	}

	return 0
}
