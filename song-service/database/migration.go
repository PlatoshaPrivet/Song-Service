package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	connString := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal("Error on connecting database: ", err)
	}

	DB = db

	// Migrating DB
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Error on migrating database: ", err)
	}

	//New migrate instance
	m, err := migrate.NewWithDatabaseInstance("file://database/migrations", "postgres", driver)
	if err != nil {
		log.Fatal("Error on connecting database: ", err)
	}
	m.Up()
}
