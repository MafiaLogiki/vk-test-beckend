package store

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"marketplace-service/internal/model"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrInvalidUsernameOrPassword = errors.New("invalid username or password")

type PostgresUserStore struct {
	DB *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{DB: db}
}

func (s *PostgresUserStore) CreateUser(user *model.User) (int64, error) {
	query := `INSERT INTO users(username, password) VALUES($1, $2) RETURNING id`
	var id int64
	err := s.DB.QueryRow(query, user.Username, user.Password).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return 0, ErrUserAlreadyExists
		}
		return 0, err
	}
	return id, nil
}

func (s *PostgresUserStore) GetUserByCredentials(username, password string) (*model.User, error) {
	query := `SELECT id, username, password FROM users WHERE username = $1 and password = $2`
	user := &model.User{}
	err := s.DB.QueryRow(query, username, password).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidUsernameOrPassword
		}
		return nil, err
	}
	return user, nil
}
