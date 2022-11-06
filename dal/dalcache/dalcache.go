package dalcache

import (
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/log"
	"github.com/devingen/sepet-cdn/cache"
	"github.com/devingen/sepet-cdn/model"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"time"
)

// updateTicker controls the frequency of underlying data updating.
var updateTicker *time.Ticker

// DALCache implements DAL interface with map cache storage
type DALCache struct {
	context    context.Context
	logger     *logrus.Logger
	Buckets    []*model.Bucket
	FileCache  cache.IFileCache
	HTTPClient *resty.Client
	apiURL     string
}

func New(ctx context.Context, fileCache cache.IFileCache, apiURL, apiKey string, dalUpdateInterval time.Duration) (*DALCache, error) {
	logger, err := log.Of(ctx)
	if err != nil {
		return nil, err
	}

	dal := &DALCache{
		context:    ctx,
		logger:     logger,
		FileCache:  fileCache,
		HTTPClient: resty.New().SetHeader("api-key", apiKey),
		apiURL:     apiURL,
	}

	buckets, _, err := dal.fetchBucketList()
	if err != nil {
		return nil, err
	}
	dal.Buckets = buckets

	// update the data periodically
	updateTicker = time.NewTicker(dalUpdateInterval)
	go func() {
		for range updateTicker.C {
			dal.Refresh()
		}
	}()

	return dal, nil
}

func (dal *DALCache) GetBucket(domain string) *model.Bucket {
	for _, bucket := range dal.Buckets {
		if core.StringValue(bucket.Domain) == domain {
			return bucket
		}
	}
	return nil
}

func (dal *DALCache) Refresh() {
	dal.logger.Info("refreshing-cache")

	buckets, _, err := dal.fetchBucketList()
	if err != nil {
		dal.logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("refreshing-cache-failed")
		return
	}
	dal.Buckets = buckets

	dal.FileCache.Invalidate(buckets)
	return
}

func (dal *DALCache) fetchBucketList() ([]*model.Bucket, *resty.Response, error) {
	var response GetBucketListResponse
	resp, err := dal.HTTPClient.R().
		SetResult(&response).
		Get(dal.apiURL + "/buckets")

	dal.logger.WithFields(logrus.Fields{
		"bucketCount": len(response.Results),
		"status":      resp.Status(),
	}).Info("retrieved-bucket-list")

	return response.Results, resp, err
}
