package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{userService: *userService}
}

func (handler *UserHandler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return router
}
