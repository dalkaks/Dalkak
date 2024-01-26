package app

import (
	"dalkak/domain/user"
	"dalkak/pkg/interfaces"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *APP) NewRouter(userService interfaces.UserService) *chi.Mux {
	router := chi.NewRouter()

	userHandler := user.NewUserHandler(userService, app.verifyMetaMaskSignature)

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(app.enableCORS)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(app.Origin))
	})

	router.Mount("/user", userHandler.Routes())

	return router
}
