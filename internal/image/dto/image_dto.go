package dto

import "mime/multipart"

type ImageUploadRequest struct {
	FileHeader *multipart.FileHeader `form:"file" binding:"required"`
}

type ImageUploadResponse struct {
	ImageUrl string `json:"imageUrl"`
}
