package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"marketplace-service/internal/config"
	"marketplace-service/internal/logger"
)

var database *sql.DB

func ConnectToDatabase(cfg *config.Config, l logger.Logger) (error) {
    databaseInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
        cfg.Postgres.Host,
        cfg.Postgres.Port,
        cfg.Postgres.User,
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

func CreateNewUser(username, password string) error {
	query := fmt.Sprintf(`
		INSERT INTO users(username, password) VALUE("%v", "%v")`, username, password)
	_, err := database.Exec(query)
	return err 
}
