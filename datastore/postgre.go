package datastore

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Postgre is our global variable for posgres datastore
var Postgre *sqlx.DB

// ConnectPostgre This function connects to postgres datastore
func ConnectPostgre() {
	var err error
	Postgre, err = sqlx.Open("postgres", "user=postgres dbname=fileServer sslmode=disable")
	if err != nil {
		panic(err)
	}

	Postgre.SetMaxIdleConns(1)
	Postgre.SetMaxOpenConns(8)

}
