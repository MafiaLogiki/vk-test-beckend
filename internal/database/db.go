package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"marketplace-service/internal/config"
	"marketplace-service/internal/logger"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrInvalidUsernameOrPassword = errors.New("invalid username or password")

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

func CreateNewUser(username, password string) (int64, error) {
	query := `INSERT INTO users(username, password) VALUES($1, $2)`
	result, err := database.Exec(query, username, password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, ErrUserAlreadyExists
		}
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

func CheckIfUserValid(username, password string) (int64, error) {
	query := `SELECT id, username FROM users WHERE username = $1 and password = $2`
	result, err := database.Query(query, username, password)
	if err != nil {
		return 0, err
	}

	defer result.Close()
	var data struct {
		Id       int64
		Username string
	}
	
	result.Next()
	if result.Scan(&data.Id, &data.Username) != nil {
		return 0, ErrInvalidUsernameOrPassword
	}

	return data.Id, nil
}
