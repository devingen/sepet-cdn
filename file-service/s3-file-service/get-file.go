package s3fs

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	fs "github.com/devingen/sepet-cdn/file-service"
	"io/ioutil"
)

// GetFile implements IFileService interface
func (s3Service S3Service) GetFile(ctx context.Context, filePath string) (*s3.GetObjectOutput, []byte, error) {
	sess := session.New(s3Service.Config)
	s3Client := s3.New(sess, s3Service.Config)

	// try to get the file
	fileMeta, err := s3Client.GetObject(&s3.GetObjectInput{Bucket: aws.String(s3Service.Bucket), Key: aws.String(filePath)})
	if err != nil {
		return nil, nil, fs.ErrorFileNotFound
	}

	fileContent, err := ioutil.ReadAll(fileMeta.Body)
	if err != nil {
		return nil, nil, err
	}

	return fileMeta, fileContent, nil
}
