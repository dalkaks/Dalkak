package user

import (
	"dalkak/pkg/interfaces"
	"dalkak/pkg/payloads"
	"dalkak/pkg/utils/httputils"
	"dalkak/pkg/utils/reflectutils"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService             interfaces.UserService
	verifyMetaMaskSignature func(next http.Handler) http.Handler
}

func NewUserHandler(userService interfaces.UserService, verifyMetaMaskSignature func(next http.Handler) http.Handler) *UserHandler {
	return &UserHandler{userService: userService, verifyMetaMaskSignature: verifyMetaMaskSignature}
}

func (handler *UserHandler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	router.With(handler.verifyMetaMaskSignature).Post("/auth", handler.authAndSignUp)

	router.Post("/refresh", handler.reissueRefresh)

	return router
}

func (handler *UserHandler) authAndSignUp(w http.ResponseWriter, r *http.Request) {
	var req payloads.UserAuthAndSignUpRequest
	err := reflectutils.GetRequestData(r, &req)
	if err != nil {
		httputils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	authTokens, tokenTime, err := handler.userService.AuthAndSignUp(req.WalletAddress, req.Signature)
	if err != nil {
		httputils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	mode := handler.userService.GetMode()
	domain := handler.userService.GetDomain()
	httputils.SetCookieRefresh(w, mode, authTokens.RefreshToken, tokenTime, domain)

	result := &payloads.UserAccessTokenResponse{
		AccessToken: authTokens.AccessToken,
	}
	if err := httputils.WriteJSON(w, http.StatusOK, result); err != nil {
		httputils.ErrorJSON(w, err, http.StatusInternalServerError)
	}
}

func (handler *UserHandler) reissueRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := httputils.GetCookieRefresh(r)
	if refreshToken == "" {
    httputils.ErrorJSON(w, errors.New("refresh token not found"), http.StatusBadRequest)
    return
  }

	authTokens, tokenTime, err := handler.userService.ReissueRefresh(refreshToken)
	if err != nil {
	  httputils.ErrorJSON(w, err, http.StatusInternalServerError)
	  return
	}

	mode := handler.userService.GetMode()
	domain := handler.userService.GetDomain()
	httputils.SetCookieRefresh(w, mode, authTokens.RefreshToken, tokenTime, domain)

	result := &payloads.UserAccessTokenResponse{
	  AccessToken: authTokens.AccessToken,
	}
	if err := httputils.WriteJSON(w, http.StatusOK, result); err != nil {
	  httputils.ErrorJSON(w, err, http.StatusInternalServerError)
	}
}
