package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"

	"marketplace-service/internal/config"
)

var database *sql.DB

func ConnectToDatabase(cfg *config.Config) (error) {
    databaseInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
        cfg.Postgres.Host,
        cfg.Postgres.Port,
        cfg.Postgres.HostName,
        cfg.Postgres.Password,
        cfg.Postgres.DBName,
    )
    
    var err error
    database, err = sql.Open("postgres", databaseInfo)
    
    if database.Ping() != nil {
        return database.Ping()
    }

    return err
}
