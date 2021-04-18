package srvcont

import (
	"bytes"
	"context"
	core "github.com/devingen/api-core"
	"github.com/devingen/api-core/log"
	"github.com/devingen/sepet-cdn/cache"
	"github.com/devingen/sepet-cdn/controller"
	"github.com/devingen/sepet-cdn/dal"
	fs "github.com/devingen/sepet-cdn/file-service"
	"github.com/devingen/sepet-cdn/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strings"
	"time"
)

// ServiceController implements IServiceController interface by using IDamgaService
type ServiceController struct {
	logger      *logrus.Logger
	FileCache   cache.IFileCache
	FileService fs.IFileService
	DAL         dal.DAL
}

// New generates new ServiceController
func New(ctx context.Context, dal dal.DAL, cache cache.IFileCache, fileService fs.IFileService) (controller.IServiceController, error) {
	logger, err := log.Of(ctx)
	if err != nil {
		return nil, err
	}

	return ServiceController{
		DAL:         dal,
		FileCache:   cache,
		FileService: fileService,
		logger:      logger,
	}, nil
}

func (sc ServiceController) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	bucketDomain := GetBucketDomainNameFromHost(r.Host)
	bucket := sc.DAL.GetBucket(bucketDomain)
	if bucket == nil {
		http.Error(w, "bucket-not-found", http.StatusNotFound)
		return
	}

	if core.StringValue(bucket.Status) != "active" {
		http.Error(w, "bucket-not-active", http.StatusGone)
		return
	}

	filePath, errorFilePath := getFilePath(bucket, r.URL.Path)

	startTime := time.Now()
	logger := sc.logger.WithFields(logrus.Fields{
		"domain":  core.StringValue(bucket.Domain),
		"folder":  core.StringValue(bucket.Folder),
		"version": core.StringValue(bucket.Version),
		"file":    filePath,
	})

	if core.BoolValue(bucket.IsCacheEnabled) {
		fileContent, fileMeta, hasCache := sc.FileCache.GetFile(filePath)
		if hasCache {
			logElapsedTime(logger, startTime, true, fileMeta.ContentLength)

			setCorsHeadersForOrigin(w, r.Header.Get("Origin"), bucket)
			http.ServeContent(w, r, filePath, fileMeta.LastModified.UTC(), bytes.NewReader(fileContent))
			return
		}
	}

	// try to get the file
	fileMeta, fileContent, err := sc.FileService.GetFile(ctx, filePath)
	if err != nil {
		if err == fs.ErrorFileNotFound {
			logger.WithFields(logrus.Fields{
				"file": filePath,
			}).Debug("file-not-found")

			// try to get the error file from cache
			if core.BoolValue(bucket.IsCacheEnabled) {
				fileContent, fileMeta, hasCache := sc.FileCache.GetFile(errorFilePath)
				if hasCache {
					logElapsedTime(logger, startTime, true, fileMeta.ContentLength)

					setCorsHeadersForOrigin(w, r.Header.Get("Origin"), bucket)
					http.ServeContent(w, r, filePath, fileMeta.LastModified.UTC(), bytes.NewReader(fileContent))
					return
				}
			}

			// try to get the error file from file server
			fileMeta, fileContent, err = sc.FileService.GetFile(ctx, errorFilePath)
			if err != nil {
				if err == fs.ErrorFileNotFound {
					logger.WithFields(logrus.Fields{
						"file": errorFilePath,
					}).Debug("error-file-not-found")

					http.Error(w, "file-not-found", http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// save the error file into cache
			sc.FileCache.SaveFile(errorFilePath, fileMeta, fileContent)
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		// save the file into cache
		sc.FileCache.SaveFile(filePath, fileMeta, fileContent)
	}

	logElapsedTime(logger, startTime, false, fileMeta.ContentLength)

	setCorsHeadersForOrigin(w, r.Header.Get("Origin"), bucket)
	http.ServeContent(w, r, filePath, fileMeta.LastModified.UTC(), bytes.NewReader(fileContent))
}

// GetBucketDomainNameFromHost returns the first subdomain
// Returns "acme" for "acme.sepet.devingen.io"
func GetBucketDomainNameFromHost(host string) string {
	dotIndex := strings.IndexByte(host, '.')
	if dotIndex < 0 {
		return host
	}
	return host[:dotIndex]
}

func getFilePath(bucket *model.Bucket, path string) (string, string) {
	if path == "/" {
		path = "/" + core.StringValue(bucket.IndexPagePath)
	}

	if core.StringValue(bucket.VersionIdentifier) == "path" {
		// version info is in the request path, no need to add the version to the file path
		filePath := core.StringValue(bucket.Folder) + path

		// get the version from the path
		parts := strings.Split(path, "/")
		version := parts[1]
		errorFilePath := core.StringValue(bucket.Folder) + "/" + version + "/" + core.StringValue(bucket.ErrorPagePath)

		return filePath, errorFilePath
	}

	filePath := core.StringValue(bucket.Folder) + "/" + core.StringValue(bucket.Version) + path
	errorFilePath := core.StringValue(bucket.Folder) + "/" + core.StringValue(bucket.Version) + "/" + core.StringValue(bucket.ErrorPagePath)
	return filePath, errorFilePath
}

func logElapsedTime(logger *logrus.Entry, startTime time.Time, fromCache bool, fileSize *int64) {

	// calculate elapsed time
	elapsedMilliseconds := int64(time.Since(startTime) / time.Millisecond)

	logger.WithFields(logrus.Fields{
		"response-time": elapsedMilliseconds,
		"from-cache":    fromCache,
		"size":          *fileSize,
	}).Debug("served")
}

func setCorsHeadersForOrigin(w http.ResponseWriter, origin string, bucket *model.Bucket) {
	if bucket.CORSConfigs == nil {
		return
	}
	corsConfigs := *bucket.CORSConfigs

	hasConfigForOrigin := false
	hasConfigForAllOrigins := false

	var configForOrigin model.CORSConfig
	var configForAllOrigins model.CORSConfig
	for _, corsConfig := range corsConfigs {
		if contains(*corsConfig.AllowedOrigins, origin) {
			configForOrigin = corsConfig
			hasConfigForOrigin = true
		}

		if contains(*corsConfig.AllowedOrigins, "*") {
			configForAllOrigins = corsConfig
			hasConfigForAllOrigins = true
		}
	}

	if hasConfigForOrigin {
		setCorsHeaders(w, configForOrigin)
	} else if hasConfigForAllOrigins {
		setCorsHeaders(w, configForAllOrigins)
	}
}

func setCorsHeaders(w http.ResponseWriter, config model.CORSConfig) {
	if config.AllowedOrigins != nil {
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(*config.AllowedOrigins, ","))
	}
	if config.AllowedMethods != nil {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(*config.AllowedMethods, ","))
	}
	if config.AllowedHeaders != nil {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(*config.AllowedHeaders, ","))
	}
	if config.ExposeHeaders != nil {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(*config.ExposeHeaders, ","))
	}
	if config.MaxAgeSeconds != nil {
		w.Header().Set("Access-Control-Max-Age", *config.MaxAgeSeconds)
	}
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}
