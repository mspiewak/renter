package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/mspiewak/renter/db"
)

// App keeps application dependencies
type App struct {
	DB                *sqlx.DB
	billRepository    *db.BillRepository
	tenantRepository  *db.TenantRepository
	paymentRepository *db.PaymentRepository
}

func main() {
	var app App

	for {
		log.Println("connecting")
		if err := app.Initialize(); err != nil {
			log.Println(err)
			time.Sleep(time.Second * 10)
			continue
		}

		break
	}
	defer app.DB.Close()

	r := mux.NewRouter()
	r.Handle("/tenant", errorHandler(app.getTenantsHandler)).Methods(http.MethodGet)
	r.Handle("/tenant/{id:[0-9a-z]+}/bill", errorHandler(app.getTenantBills)).Methods(http.MethodGet)
	r.Handle("/tenant/{id:[0-9a-z]+}/payment", errorHandler(app.getTenantPayments)).Methods(http.MethodGet)
	r.Handle("/bill", errorHandler(app.getBills)).Methods(http.MethodGet)
	r.Handle("/bill", errorHandler(app.postBill)).Methods(http.MethodPost)
	r.Handle("/bill", errorHandler(app.optionsBills)).Methods(http.MethodOptions)
	r.Handle("/stats/income", errorHandler(app.getIncome)).Methods(http.MethodGet)
	r.Handle("/cron/rent", errorHandler(app.postRent)).Methods(http.MethodPost)

	log.Println("listening")
	log.Fatal(http.ListenAndServe(":8090", commonHeaders(r)))
}

// Initialize the application
func (a *App) Initialize() error {
	dbc, err := sqlx.Connect("mysql", "root:Jg2FXug3rg@tcp(db:3306)/renter?parseTime=true")
	if err != nil {
		return fmt.Errorf("cannot connect to db server: %v", err)
	}

	a.DB = dbc
	a.billRepository = db.NewBillRepository(dbc)
	a.tenantRepository = db.NewTenantRepository(dbc)
	a.paymentRepository = db.NewPaymentRepository(dbc)
	return nil
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
	case "j2i9kzr":
		return 5
	case "hfb8yf5":
		return 6
	case "jaf8eg3":
		return 7
	}

	return 0
}
