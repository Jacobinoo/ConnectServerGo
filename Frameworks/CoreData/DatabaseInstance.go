package CoreData

import (
	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scylladb/gocqlx/v3"
)

var UserServicesDatabaseInstance *pgxpool.Pool

var StorageServicesDatabaseCluster *gocql.ClusterConfig
var StorageServicesDatabaseSession *gocqlx.Session
