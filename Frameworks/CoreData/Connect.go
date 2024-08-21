package CoreData

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func Connect() {
	db, err := pgx.Connect(context.Background(), os.Getenv("DB_URI"))
	if err != nil {
		log.Fatal("Invalid DB config:", err)
	}
	if err = db.Ping(context.Background()); err != nil {
		log.Fatal("DB unreachable:", err)
	}
	defer log.Println("Successfully connected to database and pinged it")
	DatabaseInstance = db
}
