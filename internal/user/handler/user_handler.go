package handler

import (
	"ps-beli-mang/internal/user/dto"
	"ps-beli-mang/internal/user/model"
	"ps-beli-mang/internal/user/service"
	"ps-beli-mang/pkg/helper"
	"ps-beli-mang/pkg/httphelper/response"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterAdmin(ctx echo.Context) *response.WebResponse {
	return h.Register(ctx, model.ADMIN)
}

func (h *UserHandler) RegisterUser(ctx echo.Context) *response.WebResponse {
	return h.Register(ctx, model.USER)
}

func (h *UserHandler) Register(ctx echo.Context, role model.Role) *response.WebResponse {
	var request = new(dto.RegisterRequest)
	err := ctx.Bind(&request)

	err = dto.ValidateRegisterRequest(request)
	helper.Panic400IfError(err)

	request.Role = string(role)
	result, err := h.userService.Register(ctx.Request().Context(), request)
	helper.PanicIfError(err, "error when register")

	return &response.WebResponse{
		Status:  201,
		Message: "registered successfully",
		Token:   result.AccessToken,
	}
}

func (h *UserHandler) LoginAdmin(ctx echo.Context) *response.WebResponse {
	return h.Login(ctx, model.ADMIN)
}

func (h *UserHandler) LoginUser(ctx echo.Context) *response.WebResponse {
	return h.Login(ctx, model.USER)
}

func (h *UserHandler) Login(ctx echo.Context, role model.Role) *response.WebResponse {
	var request = new(dto.LoginRequest)
	err := ctx.Bind(&request)

	request.Role = role
	err = dto.ValidateLoginRequest(request)
	helper.Panic400IfError(err)

	result, err := h.userService.Login(ctx.Request().Context(), request)
	helper.PanicIfError(err, "failed to login")

	return &response.WebResponse{
		Status:  200,
		Message: "User logged successfully",
		RawData: result,
	}

}
