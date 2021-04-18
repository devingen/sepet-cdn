package server

import (
	"context"
	"github.com/devingen/api-core/log"
	"github.com/devingen/sepet-cdn/cache/filemapcache"
	"github.com/devingen/sepet-cdn/config"
	srvcont "github.com/devingen/sepet-cdn/controller/service-controller"
	"github.com/devingen/sepet-cdn/dal/dalcache"
	s3fs "github.com/devingen/sepet-cdn/file-service/s3-file-service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmhttp"
	"net/http"
)

// New creates a new HTTP server
func New(appConfig config.App) *http.Server {

	ctx := context.Background()

	ctx, logger := getLogContext(ctx, appConfig.LogLevel)

	logger.WithFields(logrus.Fields{
		"port": appConfig.Port,
	}).Info("running-server")

	srv := &http.Server{Addr: ":" + appConfig.Port}

	fileCache, err := filemapcache.New(ctx, appConfig.CacheResetInterval)
	if err != nil {
		logger.Fatal(err)
	}

	dal, err := dalcache.New(ctx, fileCache, appConfig.ApiURL, appConfig.DalUpdateInterval)
	if err != nil {
		logger.Fatal(err)
	}

	fileService := s3fs.New(appConfig.S3)
	serviceController, err := srvcont.New(ctx, dal, fileCache, fileService)
	if err != nil {
		logger.Fatal(err)
	}

	router := mux.NewRouter()
	wrappedHandler := apmhttp.Wrap(http.HandlerFunc(serviceController.GetFile))
	router.HandleFunc("/{filePath}", wrappedHandler.ServeHTTP).Methods(http.MethodGet)

	http.HandleFunc("/", serviceController.GetFile)
	return srv
}

func getLogContext(ctx context.Context, level string) (context.Context, *logrus.Logger) {
	// create logger
	logger := logrus.New().WithFields(logrus.Fields{
		"app": "sepet-cdn",
	}).Logger

	logrusLevel, parseLevelErr := logrus.ParseLevel(level)
	if parseLevelErr == nil {
		logger.SetLevel(logrusLevel)
	}

	// add logger to the context
	ctxWithLogger := log.WithLogger(ctx, logger)
	return ctxWithLogger, logger
}
