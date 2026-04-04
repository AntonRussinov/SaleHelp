package api

import (
	"mime/multipart"
	"time"
)

type (
	HTTPError struct {
		Code             int               `json:"code"             example:"400"`
		Message          string            `json:"message"          example:"invalid body"`
		Details          string            `json:"details"          example:""`
		Timestamp        time.Time         `json:"timestamp"        example:"2024-01-15T10:00:00Z"`
		ValidationErrors []ValidationError `json:"validationErrors"`
	}

	ValidationError struct {
		Field   string `json:"field"   example:"email"`
		Message string `json:"message" example:"required"`
		Details string `json:"details" example:""`
	}

	UploadFileRequest struct {
		File *multipart.FileHeader `form:"file" swaggerignore:"true"`
	}

	// UploadFileResponse domain.FileMetadata
)
