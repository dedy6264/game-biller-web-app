package connections

import (
	"database/sql"
	"fmt"
	"gamebiller/configs"

	_ "github.com/lib/pq"
)

var (
	Db *sql.DB
)

func Connect() {
	var err error
	Db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		configs.DB_USER, configs.DB_PASS, configs.DB_NAME, configs.DB_HOST, configs.DB_PORT, configs.SSL_MODE))
	if err != nil {
		panic(err)
	}
	err = Db.Ping()
	if err != nil {
		panic(err)
	}
}

func DBconn() *sql.DB {
	return Db
}
