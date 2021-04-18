package cache

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/devingen/sepet-cdn/model"
)

// IFileCache defines the functionality of the file cache
type IFileCache interface {
	GetFile(path string) ([]byte, *s3.GetObjectOutput, bool)
	SaveFile(path string, data *s3.GetObjectOutput, buff []byte)
	Invalidate(buckets []*model.Bucket)
}
