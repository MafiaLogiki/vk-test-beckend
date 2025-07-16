package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"

	"marketplace-service/internal/config"
	"marketplace-service/internal/logger"
	"marketplace-service/internal/store"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrInvalidUsernameOrPassword = errors.New("invalid username or password")

var database *sql.DB

func ConnectToDatabase(cfg *config.Config, l logger.Logger) (store.UserStore, error) {
    databaseInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
        cfg.Postgres.Host,
        cfg.Postgres.Port,
        cfg.Postgres.User,
        cfg.Postgres.Password,
        cfg.Postgres.DBName,
    )
    
    var err error
    database, err = sql.Open("postgres", databaseInfo)

	newStore := store.NewPostgresUserStore(database)
    
    if database.Ping() != nil {
        return nil, database.Ping()
    }

    return newStore, err
}
