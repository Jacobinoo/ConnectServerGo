package CoreData

import (
	"database/sql"
	"log"
	"os"
)

func Connect() {
	db, err := sql.Open("mysql", os.Getenv("DB_URI"))
	if err != nil {
		log.Fatal("Invalid DB config:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}
	log.Println("Successfully connected to database and pinged it")
	DatabaseInstance = db
}
