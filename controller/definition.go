package controller

import (
	"net/http"
)

// IServiceController defines the functionality of the service controller
type IServiceController interface {
	GetFile(w http.ResponseWriter, r *http.Request)
}
