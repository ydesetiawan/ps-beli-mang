package dto

import (
	"github.com/go-playground/validator/v10"
	"ps-beli-mang/internal/user/model"
	"ps-beli-mang/pkg/helper"
)

type LoginRequest struct {
	Username string     `json:"username" validate:"required,min=5,max=30"`
	Password string     `json:"password" validate:"required,min=5,max=30"`
	Role     model.Role `json:"-"`
}

func ValidateLoginRequest(request *LoginRequest) error {
	validate := validator.New()
	return validate.Struct(request)
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Password string `json:"password" validate:"required,min=5,max=30"`
	Email    string `json:"email" validate:"required,email,min=1,max=255"`
	Role     string `json:"-"`
}

func NewUser(request *RegisterRequest) *model.User {
	return &model.User{
		ID:       helper.GenerateULID(),
		Username: request.Username,
		Email:    request.Email,
		Role:     request.Role,
	}

}

func ValidateRegisterRequest(request *RegisterRequest) error {
	validate := validator.New()
	return validate.Struct(request)
}

type UserResponse struct {
	AccessToken string `json:"token"`
}
