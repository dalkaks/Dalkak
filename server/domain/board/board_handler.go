package board

import (
	"dalkak/pkg/interfaces"

	"github.com/go-chi/chi/v5"
)

type BoardHandlerImpl struct {
	boardService interfaces.BoardService
}

func NewBoardHandler(boardService interfaces.BoardService) *BoardHandlerImpl {
	return &BoardHandlerImpl{boardService}
}

func (handler *BoardHandlerImpl) Routes() chi.Router {
	router := chi.NewRouter()

	return router
}
