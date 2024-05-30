package service

import (
	"golang.org/x/net/context"
	"ps-beli-mang/internal/user/dto"
	"ps-beli-mang/internal/user/model"
)

type UserService interface {
	GetUserByIdAndRole(ctx context.Context, id string, role string) (model.User, error)
	Register(context.Context, *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(context.Context, *dto.LoginRequest) (*dto.UserResponse, error)
}
