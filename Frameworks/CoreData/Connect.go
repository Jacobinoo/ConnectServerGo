package CoreData

import (
	"database/sql"
	"log"
)

func Connect() {
	db, err := sql.Open("mysql", "root:SMJyflsGzCIWuqmGSPtmcHFZxCLxQAsX@tcp(roundhouse.proxy.rlwy.net:11308)/railway")
	if err != nil {
		log.Fatal("Invalid DB config:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}
	log.Println("Successfully connected to database and pinged it")
	DatabaseInstance = db
}
