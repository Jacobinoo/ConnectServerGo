module ConnectServer

go 1.22.3

require github.com/jackc/pgx/v5 v5.6.0

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

require (
	github.com/gocql/gocql v1.6.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/scylladb/gocqlx/v3 v3.0.0
	golang.org/x/crypto v0.23.0
)

replace github.com/gocql/gocql => github.com/scylladb/gocql v1.14.3
