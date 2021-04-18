package main

import (
	"github.com/devingen/sepet-cdn/config"
	"github.com/devingen/sepet-cdn/server"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

func main() {

	var appConfig config.App
	err := envconfig.Process("sepet", &appConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	srv := server.New(appConfig)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Listen and serve failed %s", err.Error())
	}
}
