package http

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/ar-agahian/ice-assignment/internal/usecase"
	apperrors "github.com/ar-agahian/ice-assignment/pkg/errors"
	"github.com/gin-gonic/gin"
)

// FileHandler handles file upload HTTP requests
type FileHandler struct {
	fileUseCase *usecase.FileUseCase
}

// NewFileHandler creates a new FileHandler
func NewFileHandler(fileUseCase *usecase.FileUseCase) *FileHandler {
	return &FileHandler{
		fileUseCase: fileUseCase,
	}
}

// UploadFileResponse represents the response for file upload
type UploadFileResponse struct {
	FileID string `json:"fileId"`
}

// UploadFile handles POST /upload requests
func (h *FileHandler) UploadFile(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		appErr := apperrors.NewAppError("FILE_REQUIRED", "file is required", http.StatusBadRequest, nil)
		c.Error(appErr)
		return
	}

	src, err := header.Open()
	if err != nil {
		appErr := apperrors.NewAppError("INVALID_INPUT", "invalid input", http.StatusBadRequest, nil)
		c.Error(appErr)
		return
	}
	defer func() {
		if err := src.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, src); err != nil {
		appErr := apperrors.NewAppError("INTERNAL_ERROR", "internal server error", http.StatusInternalServerError, err)
		c.Error(appErr)
		return
	}

	contentType := http.DetectContentType(buf.Bytes())
	fileSize := header.Size
	filename := header.Filename

	fileID, err := h.fileUseCase.UploadFile(c.Request.Context(), usecase.UploadFileRequest{
		File:        bytes.NewReader(buf.Bytes()),
		ContentType: contentType,
		Size:        fileSize,
		Filename:    filename,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, UploadFileResponse{
		FileID: fileID,
	})
}

// RegisterRoutes registers file routes
func (h *FileHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/asset", h.UploadFile)
}
