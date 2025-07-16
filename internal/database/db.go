package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	"marketplace-service/internal/config"
	"marketplace-service/internal/logger"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrInvalidUsernameOrPassword = errors.New("invalid username or password")


func ConnectToDatabase(cfg *config.Config, l logger.Logger) (*sql.DB, error) {
    databaseInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
        cfg.Postgres.Host,
        cfg.Postgres.Port,
        cfg.Postgres.User,
        cfg.Postgres.Password,
        cfg.Postgres.DBName,
    )
    
	database, err := sql.Open("postgres", databaseInfo)
    
    if database.Ping() != nil {
        return nil, database.Ping()
    }

    return database, err
}

func CloseConnection(db *sql.DB) {
	db.Close()
}
