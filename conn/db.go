package conn

import (
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := sqlx.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatal("Erorr in connecting in database with :", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(150)
	db.SetMaxIdleConns(120)
	// gorm sql init
	DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Erorr in connecting in database with :", err)
	}
	log.Println("DB: connected!")
}

var DB2 *gorm.DB

func SecondConnectDB() {
	db, err := sqlx.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatal("Erorr in connecting in database with :", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(150)
	db.SetMaxIdleConns(120)
	// gorm sql init
	DB2, err = gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("Erorr in connecting in database with :", err)
	}
	log.Println("DB2: connected!")
}
