package app

import (
	"dalkak/domain/board"
	"dalkak/domain/user"
	"dalkak/pkg/interfaces"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *APP) NewRouter(userService interfaces.UserService, boardService interfaces.BoardService) *chi.Mux {
	router := chi.NewRouter()

	userHandler := user.NewUserHandler(userService, app.verifyMetaMaskSignature)
	boardHandler := board.NewBoardHandler(boardService)

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(app.enableCORS)
	router.Use(app.getTokenFromHeader)
	router.Use(app.processData)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(app.Origin))
	})

	router.Mount("/user", userHandler.Routes())
	router.Mount("/board", boardHandler.Routes())

	return router
}
