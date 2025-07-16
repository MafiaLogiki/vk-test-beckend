package store

import "marketplace-service/internal/model"

type UserStore interface {
	CreateUser(user *model.User) (int64, error)
	GetUserByCredentials(username, password string) (*model.User, error)
}
