package dal

import (
	"github.com/devingen/sepet-cdn/model"
)

// DAL defines the Data Access Layer for buckets
type DAL interface {
	GetBucket(domain string) *model.Bucket
	Refresh()
}
