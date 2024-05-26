package repository

import (
	"golang.org/x/net/context"
	"ps-beli-mang/internal/user/model"
)

type UserRepository interface {
	GetUserByIDAndRole(ctx context.Context, id string, role string) (model.User, error)
	GetUserByUsernameAndRole(ctx context.Context, nip string, role string) (model.User, error)
	Register(ctx context.Context, user *model.User) (string, error)
}
