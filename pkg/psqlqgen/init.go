package mysqlqgen

import (
	"fmt"
	"ps-beli-mang/configs"
	"time"

	"golang.org/x/exp/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Init(dbConfig *configs.MainConfig) *sqlx.DB {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Errors")
			fmt.Println("Recovered from panic:", r)
		}
	}()

	// Construct the connection string
	dsnString := dbConfig.GetDsnString()

	// Connect to PostgreSQL database with tracing
	db, err := sqlx.Open("postgres", dsnString)
	if err != nil {
		msg := fmt.Sprintf("Cannot connect to PostgreSQL: %s, %v", dsnString, err)
		slog.Error(msg)
		panic(msg)
	}

	// Set database connection pool settings
	db.SetMaxOpenConns(300)
	db.SetMaxIdleConns(300)
	db.SetConnMaxLifetime(3 * time.Minute)

	return db
}
