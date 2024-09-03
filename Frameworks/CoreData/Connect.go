package CoreData

import (
	"context"
	"log"
	"os"

	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scylladb/gocqlx/v3"
)

func ConnectUserServices() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("USER_SERVICES_DB_URI"))
	if err != nil {
		log.Fatal("Invalid DB config: ", err)
	}
	if err = dbpool.Ping(context.Background()); err != nil {
		log.Fatal("DB unreachable: ", err)
	}
	defer log.Println("User Services Database Available, Connected successfully!")
	UserServicesDatabaseInstance = dbpool
}

func ConnectStorageServices() {
	StorageServicesDatabaseCluster = gocql.NewCluster(os.Getenv("STORAGE_SERVICES_DB_URI"))
	StorageServicesDatabaseCluster.Keyspace = os.Getenv("STORAGE_SERVICES_DB_SHARED_SPACE")

	session, err := gocqlx.WrapSession(StorageServicesDatabaseCluster.CreateSession())
	if err != nil {
		log.Fatal(err)
	}
	defer log.Print("Storage Services Database Available, Connected successfully!")

	StorageServicesDatabaseSession = &session
}
