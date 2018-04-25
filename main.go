package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// App keeps application dependencies
type App struct {
	DB *sqlx.DB
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
	r.Handle("/", errorHandler(app.getBills)).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8090", commonHeaders(r)))
}

// Initialize the application
func (a *App) Initialize() error {
	var err error
	a.DB, err = sqlx.Connect("mysql", "root:root@tcp(localhost:3306)/renter?parseTime=true")
	return err
}

func (a *App) getTenantsHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return getTenants(a.DB)
}

func (a *App) getBills(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	log.Println("bills")
	return getBills(a.DB)
}

func (a *App) getTenantBills(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	tenantID := getRealTenantID(vars["id"])

	return getTenantBills(a.DB, tenantID)
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
