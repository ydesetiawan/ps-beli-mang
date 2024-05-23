package service

import (
	"ps-beli-mang/internal/user/dto"
	"ps-beli-mang/internal/user/repository"
	"ps-beli-mang/pkg/bcrypt"
	"ps-beli-mang/pkg/errs"
	"ps-beli-mang/pkg/middleware"
)

type UserService interface {
	Register(request *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(*dto.LoginRequest) (*dto.UserResponse, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserServiceImpl(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) Register(request *dto.RegisterRequest) (*dto.UserResponse, error) {
	user := dto.NewUser(request)
	hashedPassword, _ := bcrypt.HashPassword(request.Password)
	user.Password = hashedPassword

	id, err := s.userRepository.Register(user)
	if err != nil {
		return &dto.UserResponse{}, err
	}
	token, _ := middleware.GenerateJWT(id, request.Role)
	return &dto.UserResponse{
		AccessToken: token,
	}, nil
}

func (s *userService) Login(req *dto.LoginRequest) (*dto.UserResponse, error) {
	response := &dto.UserResponse{}

	user, err := s.userRepository.GetUserByUsernameAndRole(req.Username, string(req.Role))
	if err != nil {
		return response, errs.NewErrDataNotFound("user not found ", req.Username, errs.ErrorData{})
	}
	err = bcrypt.ComparePassword(req.Password, user.Password)
	if err != nil {
		return response, errs.NewErrBadRequest("password is wrong ")
	}

	token, _ := middleware.GenerateJWT(user.ID, user.Role)

	response = &dto.UserResponse{
		AccessToken: token,
	}

	return response, nil
}
