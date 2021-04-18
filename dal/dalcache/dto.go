package dalcache

import "github.com/devingen/sepet-cdn/model"

type GetBucketListResponse struct {
	Results []*model.Bucket `json:"results"`
}
