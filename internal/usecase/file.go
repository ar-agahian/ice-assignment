package usecase

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/ar-agahian/ice-assignment/internal/interfaces/client"
	apperrors "github.com/ar-agahian/ice-assignment/pkg/errors"
)

const (
	maxFileSize = 10 * 1024 * 1024
)

var (
	allowedMimeTypes = []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"application/pdf",
		"text/plain",
	}
)

// FileUseCase handles file upload business logic
type FileUseCase struct {
	storageRepo client.IFileStorage
}

// NewFileUseCase creates a new FileUseCase
func NewFileUseCase(storageRepo client.IFileStorage) *FileUseCase {
	return &FileUseCase{
		storageRepo: storageRepo,
	}
}

// UploadFileRequest represents the request to upload a file
type UploadFileRequest struct {
	File        io.Reader
	ContentType string
	Size        int64
	Filename    string
}

// UploadFile uploads a file to S3 and returns the file ID
func (uc *FileUseCase) UploadFile(ctx context.Context, req UploadFileRequest) (string, error) {
	if req.File == nil {
		return "", apperrors.NewAppError("FILE_REQUIRED", "file is required", http.StatusBadRequest, nil)
	}
	if req.Size == 0 {
		return "", apperrors.NewAppError("FILE_EMPTY", "file cannot be empty", http.StatusBadRequest, nil)
	}
	if req.Size > maxFileSize {
		return "", apperrors.NewAppError(
			"FILE_TOO_LARGE",
			fmt.Sprintf("file size exceeds maximum allowed size of %d bytes", maxFileSize),
			http.StatusBadRequest,
			nil,
		)
	}

	if !uc.isValidContentType(req.ContentType) {
		return "", apperrors.NewAppError(
			"INVALID_FILE_TYPE",
			"file type not allowed",
			http.StatusBadRequest,
			nil,
		)
	}
	fileID, err := uc.storageRepo.Upload(ctx, req.File, req.ContentType)
	if err != nil {
		return "", err
	}
	return fileID, nil
}

// isValidContentType checks if the content type is allowed (business rule)
func (uc *FileUseCase) isValidContentType(contentType string) bool {
	if contentType == "" {
		return false
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	for _, allowed := range allowedMimeTypes {
		if mediaType == allowed {
			return true
		}
	}
	return false
}
