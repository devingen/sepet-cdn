package s3fs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/devingen/sepet-cdn/config"
)

// S3Service implements ISepetService interface with database connection
type S3Service struct {
	Bucket string
	Config *aws.Config
}

// New generates new S3Service
func New(envConfig config.S3) S3Service {

	if envConfig.Endpoint != "" {
		// configure to use MinIO Server
		return S3Service{
			Bucket: envConfig.Bucket,
			Config: &aws.Config{
				Credentials:      credentials.NewStaticCredentials(envConfig.AccessKeyID, envConfig.AccessKey, ""),
				Endpoint:         aws.String(envConfig.Endpoint),
				Region:           aws.String(envConfig.Region),
				DisableSSL:       aws.Bool(true),
				S3ForcePathStyle: aws.Bool(true),
			},
		}
	}

	return S3Service{
		Bucket: envConfig.Bucket,
		Config: &aws.Config{
			Credentials: credentials.NewStaticCredentials(envConfig.AccessKeyID, envConfig.AccessKey, ""),
			Endpoint:    aws.String(envConfig.Endpoint),
			Region:      aws.String(envConfig.Region),
		},
	}
}
