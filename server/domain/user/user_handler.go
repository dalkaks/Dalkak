package user

import (
	"dalkak/pkg/interfaces"
	"dalkak/pkg/payloads"
	"dalkak/pkg/utils/httputils"
	"dalkak/pkg/utils/parseutils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandlerImpl struct {
	userService             interfaces.UserService
	verifyMetaMaskSignature func(next http.Handler) http.Handler
}

func NewUserHandler(userService interfaces.UserService, verifyMetaMaskSignature func(next http.Handler) http.Handler) *UserHandlerImpl {
	return &UserHandlerImpl{userService: userService, verifyMetaMaskSignature: verifyMetaMaskSignature}
}

func (handler *UserHandlerImpl) Routes() chi.Router {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	router.With(handler.verifyMetaMaskSignature).Post("/auth", handler.RouteAuthAndSignUp)

	router.Post("/refresh", handler.RouteReissueRefresh)

	router.Post("/logout", handler.RouteLogout)

	router.Post("/media/presigned", handler.RouteCreatePresignedURL)

	router.Get("/media", handler.RouteGetUserMedia)

	router.Post("/media/confirm", handler.RouteConfirmMediaUpload)

	return router
}

func (handler *UserHandlerImpl) RouteAuthAndSignUp(w http.ResponseWriter, r *http.Request) {
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

func (handler *UserHandlerImpl) RouteReissueRefresh(w http.ResponseWriter, r *http.Request) {
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

func (handler *UserHandlerImpl) RouteLogout(w http.ResponseWriter, r *http.Request) {
	httputils.DeleteCookieRefresh(w)
	httputils.WriteJSONAndHandleError(w, http.StatusOK, nil, httputils.HandleAppError)
}

func (handler *UserHandlerImpl) RouteCreatePresignedURL(w http.ResponseWriter, r *http.Request) {
	userInfo, err := parseutils.GetUserInfoData(r)
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

func (handler *UserHandlerImpl) RouteGetUserMedia(w http.ResponseWriter, r *http.Request) {
	userInfo, err := parseutils.GetUserInfoData(r)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	var query payloads.UserGetMediaRequest
	err = parseutils.GetQuery(r, &query)
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

func (handler *UserHandlerImpl) RouteConfirmMediaUpload(w http.ResponseWriter, r *http.Request) {
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

func handleAppErrorAndDeleteCookieRefresh(w http.ResponseWriter, err error) {
	httputils.DeleteCookieRefresh(w)
	httputils.HandleAppError(w, err)
}
