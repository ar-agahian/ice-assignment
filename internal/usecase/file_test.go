package usecase

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	apperrors "github.com/ar-agahian/ice-assignment/pkg/errors"
	"github.com/ar-agahian/ice-assignment/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadFile(t *testing.T) {
	tests := []struct {
		name          string
		req           UploadFileRequest
		setupMocks    func(*mocks.MockIFileStorage)
		expectedError error
	}{
		{
			name: "successful upload",
			req: UploadFileRequest{
				File:        strings.NewReader("test content"),
				ContentType: "text/plain",
				Size:        12,
			},
			setupMocks: func(storage *mocks.MockIFileStorage) {
				storage.On("Upload", mock.Anything, mock.Anything, "text/plain").Return("file-123", nil)
			},
			expectedError: nil,
		},
		{
			name: "nil file",
			req: UploadFileRequest{
				File:        nil,
				ContentType: "text/plain",
				Size:        4,
			},
			setupMocks: func(storage *mocks.MockIFileStorage) {
				// No mocks needed, validation fails early
			},
			expectedError: apperrors.NewAppError("FILE_REQUIRED", "file is required", http.StatusBadRequest, nil),
		},
		{
			name: "empty file",
			req: UploadFileRequest{
				File:        strings.NewReader(""),
				ContentType: "text/plain",
				Size:        0,
			},
			setupMocks: func(storage *mocks.MockIFileStorage) {
				// No mocks needed, validation fails early
			},
			expectedError: apperrors.NewAppError("FILE_EMPTY", "file cannot be empty", http.StatusBadRequest, nil),
		},
		{
			name: "file too large",
			req: UploadFileRequest{
				File:        strings.NewReader("test"),
				ContentType: "text/plain",
				Size:        11 * 1024 * 1024, // 11MB
			},
			setupMocks: func(storage *mocks.MockIFileStorage) {
				// No mocks needed, validation fails early
			},
			expectedError: apperrors.NewAppError("FILE_TOO_LARGE", "file size exceeds maximum allowed size of 10485760 bytes", http.StatusBadRequest, nil),
		},
		{
			name: "invalid content type",
			// Use binary content that will be detected as application/octet-stream,
			// which will cause fallback to provided ContentType (which is invalid)
			req: UploadFileRequest{
				File:        bytes.NewReader([]byte{0x4D, 0x5A, 0x90, 0x00}), // MZ header (executable)
				ContentType: "application/x-msdownload",
				Size:        4,
			},
			setupMocks: func(storage *mocks.MockIFileStorage) {
				// No mocks needed, validation fails early
			},
			expectedError: apperrors.NewAppError("INVALID_FILE_TYPE", "file type not allowed", http.StatusBadRequest, nil),
		},
		{
			name: "storage error",
			req: UploadFileRequest{
				File:        strings.NewReader("test"),
				ContentType: "text/plain",
				Size:        4,
			},
			setupMocks: func(storage *mocks.MockIFileStorage) {
				storage.On("Upload", mock.Anything, mock.Anything, "text/plain").Return("", errors.New("storage error"))
			},
			expectedError: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewMockIFileStorage(t)
			tt.setupMocks(storage)

			uc := NewFileUseCase(storage)
			fileID, err := uc.UploadFile(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if appErr, ok := apperrors.AsAppError(tt.expectedError); ok {
					actualErr, ok := apperrors.AsAppError(err)
					assert.True(t, ok, "expected AppError")
					assert.Equal(t, appErr.Code, actualErr.Code)
				} else {
					assert.NotNil(t, err)
				}
				assert.Empty(t, fileID)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, fileID)
				storage.AssertExpectations(t)
			}
		})
	}
}

