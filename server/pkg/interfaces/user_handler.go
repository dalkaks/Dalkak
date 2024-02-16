package interfaces

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandler interface {
	Routes() chi.Router

	RouteAuthAndSignUp(w http.ResponseWriter, r *http.Request)
	RouteReissueRefresh(w http.ResponseWriter, r *http.Request)
	RouteLogout(w http.ResponseWriter, r *http.Request)

	RouteCreatePresignedURL(w http.ResponseWriter, r *http.Request)
	RouteGetUserMedia(w http.ResponseWriter, r *http.Request)
	RouteConfirmMediaUpload(w http.ResponseWriter, r *http.Request)
}
