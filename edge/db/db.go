package db

import (
	"log"
	"patient-edge/config"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var DB *sqlx.DB

func init() {
	var err error
	DB, err = sqlx.Open("mysql", config.Config.Edge.DataSource)
	if err != nil {
		log.Fatal(errors.Wrap(err, "open db error"))
	}
}
