package user

import (
	"dalkak/pkg/interfaces"
	"dalkak/pkg/payloads"
	"dalkak/pkg/utils/httputils"
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

	router.Post("/logout", handler.logout)

	router.Post("/media/presigned", handler.createPresignedURL)

	router.Get("/media", handler.getUserMedia)

	router.Post("/media/confirm", handler.confirmMediaUpload)

	return router
}

func (handler *UserHandler) authAndSignUp(w http.ResponseWriter, r *http.Request) {
	var req payloads.UserAuthAndSignUpRequest
	err := httputils.ReadJSON(w, r, &req)
	if err != nil {
		handleAppErrorAndDeleteCookieRefresh(w, err)
		return
	}

	authTokens, tokenTime, err := handler.userService.AuthAndSignUp(req.WalletAddress, req.Signature)
	if err != nil {
		handleAppErrorAndDeleteCookieRefresh(w, err)
		return
	}

	mode := handler.userService.GetMode()
	domain := handler.userService.GetDomain()
	err = httputils.SetCookieRefresh(w, mode, authTokens.RefreshToken, tokenTime, domain)
	if err != nil {
		handleAppErrorAndDeleteCookieRefresh(w, err)
		return
	}

	result := &payloads.UserAccessTokenResponse{
		AccessToken: authTokens.AccessToken,
	}
	if err := httputils.WriteJSON(w, http.StatusOK, result); err != nil {
		handleAppErrorAndDeleteCookieRefresh(w, err)
	}
}

func (handler *UserHandler) reissueRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := httputils.GetCookieRefresh(r)
	if err != nil {
		handleAppErrorAndDeleteCookieRefresh(w, err)
		return
	}

	authTokens, tokenTime, err := handler.userService.ReissueRefresh(refreshToken)
	if err != nil {
		handleAppErrorAndDeleteCookieRefresh(w, err)
		return
	}

	mode := handler.userService.GetMode()
	domain := handler.userService.GetDomain()
	err = httputils.SetCookieRefresh(w, mode, authTokens.RefreshToken, tokenTime, domain)
	if err != nil {
		handleAppErrorAndDeleteCookieRefresh(w, err)
		return
	}

	result := &payloads.UserAccessTokenResponse{
		AccessToken: authTokens.AccessToken,
	}
	if err := httputils.WriteJSON(w, http.StatusOK, result); err != nil {
		handleAppErrorAndDeleteCookieRefresh(w, err)
	}
}

func (handler *UserHandler) logout(w http.ResponseWriter, r *http.Request) {
	httputils.DeleteCookieRefresh(w)
	httputils.WriteJSONAndHandleError(w, http.StatusOK, nil, httputils.HandleAppError)
}

func handleAppErrorAndDeleteCookieRefresh(w http.ResponseWriter, err error) {
	httputils.DeleteCookieRefresh(w)
	httputils.HandleAppError(w, err)
}

func (handler *UserHandler) createPresignedURL(w http.ResponseWriter, r *http.Request) {
	userInfo, err := httputils.GetUserInfoData(r)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	var req payloads.UserCreateMediaRequest
	err = httputils.ReadJSON(w, r, &req)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	result, err := handler.userService.CreatePresignedURL(userInfo, &req)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	httputils.WriteJSONAndHandleError(w, http.StatusOK, result, httputils.HandleAppError)
}

func (handler *UserHandler) getUserMedia(w http.ResponseWriter, r *http.Request) {
	userInfo, err := httputils.GetUserInfoData(r)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	var query payloads.UserGetMediaRequest
	err = httputils.GetQuery(r, &query)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	result, err := handler.userService.GetUserMedia(userInfo, &query)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	httputils.WriteJSONAndHandleError(w, http.StatusOK, result, httputils.HandleAppError)
}

func (handler *UserHandler) confirmMediaUpload(w http.ResponseWriter, r *http.Request) {
	var req payloads.UserConfirmMediaRequest
	err := httputils.ReadJSON(w, r, &req)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	err = handler.userService.ConfirmMediaUpload(&req)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	httputils.WriteJSONAndHandleError(w, http.StatusOK, nil, httputils.HandleAppError)
}
