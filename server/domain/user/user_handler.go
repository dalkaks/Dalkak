package user

import (
	"dalkak/pkg/interfaces"
	"dalkak/pkg/payloads"
	"dalkak/pkg/utils/reflectutils"
	"encoding/json"
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

	return router
}

func (handler *UserHandler) authAndSignUp(w http.ResponseWriter, r *http.Request) {
	var req payloads.UserAuthAndSignUpRequest
	err := reflectutils.GetRequestData(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authTokens, err := handler.userService.AuthAndSignUp(req.WalletAddress, req.Signature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authTokens)
}
