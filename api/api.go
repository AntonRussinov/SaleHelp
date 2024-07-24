package api

import (
	"mime/multipart"
	"time"
)

type (
	HTTPError struct {
		Code             int               `json:"code"`
		Message          string            `json:"message"`
		Details          string            `json:"details"`
		Timestamp        time.Time         `json:"timestamp"`
		ValidationErrors []ValidationError `json:"validationErrors"`
	}

	ValidationError struct {
		Field   string `json:"field"`
		Message string `json:"message"`
		Details string `json:"details"`
	}

	UploadFileRequest struct {
		File *multipart.FileHeader `form:"file" swaggerignore:"true"`
	}

	// UploadFileResponse domain.FileMetadata
)
