package common

import "time"

const SignedUrlUriTemplate = "/api/internal/v1/files/%d/signed_url"

type (
	FileInformation struct {
		// The unique identifier assigned to the file by the files service.
		FileID uint64 `json:"file_id"`

		// Checksum of the file.
		Checksum string `json:"checksum,omitempty"`

		// Size of the file in storage.
		Size int64 `json:"size,omitempty"`

		// Last modification timestamps for the file.
		UpdatedAt time.Time `json:"updated_at,omitempty"`
	}

	FileContents struct {
		Size int64  `json:"size"`
		Data []byte `json:"data"`
	}

	// CommonFileInfoResponse - defines the response structure for create file, update file
	// and get file requests.
	CommonFileInfoResponse struct {
		File         FileInformation `json:"file,omitempty"`
		RequestID    string          `json:"request_id"`
		ResponseTime time.Time       `json:"response_time"`
	}

	// CommonFileContentResponse - defines the response structure for create file, update file
	// and get file requests.
	CommonFileContentsResponse struct {
		FileContents FileContents `json:"file_contents"`
		RequestID    string       `json:"request_id"`
		ResponseTime time.Time    `json:"response_time"`
	}

	// SignedUrlResponse - defines the response structure for get signed url requests.
	SignedUrlResponse struct {
		RequestID    string    `json:"request_id"`
		ResponseTime time.Time `json:"response_time"`
		FileName     string    `json:"file_name,omitempty"`
		SignedUrl    string    `json:"url,omitempty"`
	}
)
