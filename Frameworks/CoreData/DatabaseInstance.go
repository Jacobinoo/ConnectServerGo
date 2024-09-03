// CoreData
package CoreData

import (
	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/scylladb/gocqlx/v3"
)

// CoreData pointer for accessing user services database pool
var UserServicesDatabaseInstance *pgxpool.Pool

// CoreData pointer for accessing main storage services database cluster
var StorageServicesDatabaseCluster *gocql.ClusterConfig

// CoreData pointer for accessing main storage services database session
var StorageServicesDatabaseSession *gocqlx.Session
