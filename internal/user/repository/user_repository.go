package repository

import (
	"ps-beli-mang/internal/user/model"
)

type UserRepository interface {
	GetUserByIDAndRole(id string, role string) (model.User, error)
	GetUserByUsernameAndRole(nip string, role string) (model.User, error)
	Register(user *model.User) (string, error)
}
