package fs

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ErrorFileNotFound used when the file is not found
var ErrorFileNotFound = errors.New("file-not-found")

// IFileService defines the functionality of the file service
type IFileService interface {
	GetFile(ctx context.Context, filePath string) (*s3.GetObjectOutput, []byte, error)
}
