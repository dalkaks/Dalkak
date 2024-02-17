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

	router.Post("/refresh", handler.RouteReissueAccessToken)

	router.Post("/media/presigned", handler.RouteCreatePresignedURL)

	router.Get("/media", handler.RouteGetUserMedia)

	router.Post("/media/confirm", handler.RouteConfirmMediaUpload)

	router.Delete("/media", handler.RouteDeleteUserMedia)

	return router
}

func (handler *UserHandlerImpl) RouteAuthAndSignUp(w http.ResponseWriter, r *http.Request) {
	var req payloads.UserAuthAndSignUpRequest
	err := httputils.ReadJSON(w, r, &req)
	if err != nil {
	httputils.HandleAppError(w, err)
		return
	}

	result, err := handler.userService.AuthAndSignUp(&req)
	if err != nil {
	httputils.HandleAppError(w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, result); err != nil {
	httputils.HandleAppError(w, err)
	}
}

func (handler *UserHandlerImpl) RouteReissueAccessToken(w http.ResponseWriter, r *http.Request) {
	var req payloads.UserReissueAccessTokenRequest
	err := httputils.ReadJSON(w, r, &req)
	if err != nil {
	httputils.HandleAppError(w, err)
		return
	}

	result, err := handler.userService.ReissueAccessToken(&req)
	if err != nil {
	httputils.HandleAppError(w, err)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, result); err != nil {
	httputils.HandleAppError(w, err)
	}
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
	userInfo, err := parseutils.GetUserInfoData(r)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	var req payloads.UserConfirmMediaRequest
	err = httputils.ReadJSON(w, r, &req)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	err = handler.userService.ConfirmMediaUpload(userInfo, &req)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	httputils.WriteJSONAndHandleError(w, http.StatusOK, nil, httputils.HandleAppError)
}

func (handler *UserHandlerImpl) RouteDeleteUserMedia(w http.ResponseWriter, r *http.Request) {
	userInfo, err := parseutils.GetUserInfoData(r)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	var req payloads.UserDeleteMediaRequest
	err = httputils.ReadJSON(w, r, &req)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	err = handler.userService.DeleteUserMedia(userInfo, &req)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	httputils.WriteJSONAndHandleError(w, http.StatusOK, nil, httputils.HandleAppError)
}
