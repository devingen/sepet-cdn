package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// BucketStatus is the type for the Bucket's status
type BucketStatus string

const (
	// BucketStatusActive is status active
	BucketStatusActive BucketStatus = "active"
)

// Bucket defines the MongoDB and JSON structure of the bucket data
type Bucket struct {
	// DBRef fields
	Ref      string             `bson:"_ref,omitempty" json:"_ref,omitempty"`
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Database string             `bson:"_db,omitempty" json:"_db,omitempty"`

	// common model fields
	CreatedAt *time.Time `json:"_created,omitempty" bson:"_created,omitempty"`
	UpdatedAt *time.Time `json:"_updated,omitempty" bson:"_updated,omitempty"`
	Revision  int        `json:"_revision,omitempty" bson:"_revision,omitempty"`

	// Domain of the bucket used in the URL. E.g. 'acme' for 'acme.sepet.devingen.io'
	Domain *string `json:"domain,omitempty" bson:"domain,omitempty"`

	// Region of the bucket for caching the files in the proper CDN region.
	Region *string `json:"region,omitempty" bson:"region,omitempty"`

	// Folder that keeps the files. Used for preventing issues when domain name changes.
	Folder *string `json:"folder,omitempty" bson:"folder,omitempty"`

	// Version of the bucket to serve the default content. Default value is 'default'.
	//   E.g. 'acme''s folder name is a1b2c3 and the version is 0.0.1. When 'acme.sepet.devingen.io/...' is
	//   requested, it'll load the content under the 'ROOT/a1b2c3/0.0.1/...'.
	Version *string `json:"version,omitempty" bson:"version,omitempty"`

	// VersionIdentifier defines where the bucket version should be retrieved while reading the file from CDN.
	//    Should be one of 'header' or 'path'.
	//    If the version identifier is 'header';
	//      * CDN will get the version from the bucket data. Requesting
	//          'acme.sepet.devingen.io/...' would get the data from the current version.
	//      * only the current version will be cached in CDN
	//    If the version identifier is 'path'
	// 		* the path must contain the version when accessing the CND
	//      * CDN will get the version from the request path. Requesting
	//          'acme.sepet.devingen.io/0.0.1/...' will get the data from version '0.0.1'.
	//      * any version will be cached in CDN
	VersionIdentifier *string `json:"versionIdentifier,omitempty" bson:"versionIdentifier,omitempty"`

	// IndexPagePath is used for serving Single Page Web applications. This file is loaded when the root URL is called.
	IndexPagePath *string `json:"indexPagePath,omitempty" bson:"indexPagePath,omitempty"`

	// ErrorPagePath is served when the file is not found in the folder. It can be used to show a custom error page or
	//   to forward all sub routes to the index page for SPA routing.
	ErrorPagePath *string `json:"errorPagePath,omitempty" bson:"errorPagePath,omitempty"`

	// IsCacheEnabled is used by CDN to cache the file for the next request. Useful for serving static content
	//   that's fetched frequently.
	IsCacheEnabled *bool `json:"isCacheEnabled,omitempty" bson:"isCacheEnabled,omitempty"`

	// IsVersioningEnabled is used by API and CDN to allow requests with specific versions.
	IsVersioningEnabled *bool `json:"isVersioningEnabled,omitempty" bson:"isVersioningEnabled,omitempty"`

	// Status determines the bucket status. Should be one of 'active' ... (Not active status is not supported yet)
	Status *string `json:"status,omitempty" bson:"status,omitempty"`

	// CORSConfigs contains the CORS header configuration. Each config in the array
	// is supposed the be used for different rules for different origins.
	CORSConfigs *[]CORSConfig `json:"corsConfigs,omitempty" bson:"corsConfigs,omitempty"`

	// ResponseHeaders contains the headers returned to all get file responses from CDN.
	ResponseHeaders *map[string]string `json:"responseHeaders,omitempty" bson:"responseHeaders,omitempty"`
}

// AddCreationFields adds the necessary fields before inserting into database
func (b *Bucket) AddCreationFields() {
	b.ID = primitive.NewObjectID()
	now := time.Now()
	b.CreatedAt = &now
	b.UpdatedAt = &now
	b.Revision = 1
}

// PrepareUpdateFields sets the UpdatedAt and deletes the Revision. Giving 0 value to Revision results bson
// ignoring the revision field in $set function. It's incremented by the $inc command
func (b *Bucket) PrepareUpdateFields() {
	b.Revision = 0
	now := time.Now()
	b.UpdatedAt = &now
}

type CORSConfig struct {
	AllowedHeaders *[]string `json:"allowedHeaders,omitempty" bson:"allowedHeaders,omitempty"`
	AllowedMethods *[]string `json:"allowedMethods,omitempty" bson:"allowedMethods,omitempty"`
	AllowedOrigins *[]string `json:"allowedOrigins,omitempty" bson:"allowedOrigins,omitempty"`
	ExposeHeaders  *[]string `json:"exposeHeaders,omitempty" bson:"exposeHeaders,omitempty"`
	MaxAgeSeconds  *string   `json:"maxAgeSeconds,omitempty" bson:"maxAgeSeconds,omitempty"`
}
