package mysql

import (
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func ConnectDB() *sqlx.DB {
	var db *sqlx.DB
	var err error
	if os.Getenv("MAINDB") != "" {
		db, err = sqlx.Open("mysql", os.Getenv("DSN"))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db, err = sqlx.Open("mysql",
			os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME"))
		if err != nil {
			log.Fatal(err)
		}
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(150)
	db.SetMaxIdleConns(120)
	return db
}
func SecondConnectDB() *sqlx.DB {
	var db *sqlx.DB
	var err error
	if os.Getenv("SCNDDB") != "" {
		db, err = sqlx.Open("mysql", os.Getenv("SCNDDB"))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db, err = sqlx.Open("mysql",
			os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME"))
		if err != nil {
			log.Fatal(err)
		}
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(150)
	db.SetMaxIdleConns(120)
	return db
}
