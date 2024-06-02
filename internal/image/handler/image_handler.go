package handler

import (
	"ps-beli-mang/internal/image/dto"
	"ps-beli-mang/internal/image/service"
	"ps-beli-mang/internal/user/model"
	userService "ps-beli-mang/internal/user/service"
	"ps-beli-mang/pkg/base/handler"
	"ps-beli-mang/pkg/helper"
	"ps-beli-mang/pkg/httphelper/response"

	"github.com/labstack/echo/v4"
)

type ImageHandler struct {
	imageService service.ImageService
	userService  userService.UserService
}

func NewImageHandler(imageService service.ImageService, userService userService.UserService) *ImageHandler {
	return &ImageHandler{imageService: imageService, userService: userService}
}

func hasAuthorizeRoleAdmin(ctx echo.Context, h *ImageHandler) (string, error) {
	userId, err := handler.GetUserId(ctx)
	if err != nil {
		return "", err
	}

	_, err = h.userService.GetUserByIdAndRole(ctx.Request().Context(), userId, string(model.ADMIN))
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (h *ImageHandler) UploadImage(ctx echo.Context) *response.WebResponse {
	_, err := hasAuthorizeRoleAdmin(ctx, h)
	helper.PanicIfError(err, "user unauthorized")

	fileHeader, err := ctx.FormFile("file")
	helper.Panic400IfError(err)

	file, err := fileHeader.Open()
	helper.Panic400IfError(err)

	fileUrl, err := h.imageService.UploadImage(file, fileHeader)
	helper.PanicIfError(err, "failed to upload image")

	return &response.WebResponse{
		Status:  200,
		Message: "File uploaded successfully",
		Data:    dto.ImageUploadResponse{ImageUrl: fileUrl},
	}
}
