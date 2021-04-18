package filemapcache

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/log"
	"github.com/devingen/sepet-cdn/model"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

// resetTicker controls the frequency of underlying cache resetting.
var resetTicker *time.Ticker

// FileMapCache implements IFileCache interface with map cache storage
type FileMapCache struct {
	logger       *logrus.Logger
	contentCache sync.Map
	metaCache    sync.Map
}

func New(ctx context.Context, cacheResetInterval time.Duration) (*FileMapCache, error) {
	logger, err := log.Of(ctx)
	if err != nil {
		return nil, err
	}

	cache := &FileMapCache{
		logger:       logger,
		contentCache: sync.Map{},
		metaCache:    sync.Map{},
	}

	// update the data periodically
	resetTicker = time.NewTicker(cacheResetInterval)
	go func() {
		for range resetTicker.C {
			cache.Reset()
		}
	}()

	return cache, nil
}

func (mc *FileMapCache) GetFile(path string) ([]byte, *s3.GetObjectOutput, bool) {
	meta, hasMeta := mc.metaCache.Load(path)
	buff, hasBuff := mc.contentCache.Load(path)

	mc.logger.WithFields(logrus.Fields{
		"exists": hasBuff && hasMeta,
		"path":   path,
	}).Debug("getting-file-from-cache")

	if hasMeta && hasBuff {
		return buff.([]byte), meta.(*s3.GetObjectOutput), true
	}
	return nil, nil, false
}

func (mc *FileMapCache) SaveFile(path string, data *s3.GetObjectOutput, buff []byte) {
	mc.logger.WithFields(logrus.Fields{
		"path": path,
	}).Debug("saving-file-into-cache")

	mc.contentCache.Store(path, buff)
	mc.metaCache.Store(path, data)
}

func (mc *FileMapCache) Reset() {
	mc.logger.Info("resetting-cache")
	mc.contentCache = sync.Map{}
	mc.metaCache = sync.Map{}
}

func (mc *FileMapCache) Invalidate(buckets []*model.Bucket) {
	mc.logger.Info("invalidating-cache")

	pathPrefixesToKeep := map[string]bool{}
	for _, bucket := range buckets {
		if core.StringValue(bucket.Status) != "active" || !core.BoolValue(bucket.IsCacheEnabled) {
			// skip the bucket if the status is not active or caching is not enabled
			continue
		}

		if core.StringValue(bucket.VersionIdentifier) == "path" {
			// if the version identifier is path, files from different versions may have been cached.
			// we need to keep the files from all the versions of the bucket.
			// so keep all the paths starting for the bucket.
			prefix := core.StringValue(bucket.Folder) + "/"
			pathPrefixesToKeep[prefix] = true
			continue
		}

		// keep the files of the active version of the bucket
		// this will remove the cache for older version if the version is changed
		prefix := core.StringValue(bucket.Folder) + "/" + core.StringValue(bucket.Version)
		pathPrefixesToKeep[prefix] = true
	}

	mc.metaCache.Range(func(key, value interface{}) bool {
		for prefix := range pathPrefixesToKeep {
			if strings.Index(key.(string), prefix) == 0 {
				// keep the file
				return true
			}
		}
		mc.logger.WithFields(logrus.Fields{
			"path": key,
		}).Debug("removing-file-from-cache")

		mc.metaCache.Delete(key)
		mc.contentCache.Delete(key)
		return true
	})
}
