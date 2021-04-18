package config

import "time"

// App defines the environment variable configuration for the whole app
type App struct {
	// Port is the port of the HTTP server.
	Port string `envconfig:"cdn_port" default:"80"`

	// LogLevel defines the log level.
	LogLevel string `envconfig:"cdn_log_level" default:"info"`

	// DalUpdateInterval is the data refresh time interval.
	DalUpdateInterval time.Duration `envconfig:"dal_update_interval" default:"1m"`

	// CacheResetInterval is the data clean time interval.
	CacheResetInterval time.Duration `envconfig:"cache_reset_interval" default:"1h"`

	// ApiURL is the URL of the Sepet API to get buckets.
	ApiURL string `envconfig:"api_url" required:"true"`

	// S3 is the configuration of the S3 server.
	S3 S3 `envconfig:"s3"`
}

// S3 defines the environment variable configuration for AWS S3 or MinIO
type S3 struct {
	// Endpoint is the URL of the file server to connect to. If empty, the connection is made to the AWS S3 servers.
	// Used to connect a local MinIO server for development and integration tests.
	Endpoint string `envconfig:"endpoint"`

	// AccessKeyID is the AWS access key ID.
	AccessKeyID string `envconfig:"access_key_id" required:"true"`

	// AccessKey is the AWS access key.
	AccessKey string `envconfig:"secret_access_key" required:"true"`

	// Region is the AWS region.
	Region string `envconfig:"region" required:"true"`

	// Bucket is the bucket to connect.
	Bucket string `envconfig:"bucket" default:"sepet"`
}
