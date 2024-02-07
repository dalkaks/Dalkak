package board

import (
	"dalkak/pkg/interfaces"
	"dalkak/pkg/utils/httputils"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type BoardHandler struct {
	boardService interfaces.BoardService
}

func NewBoardHandler(boardService interfaces.BoardService) *BoardHandler {
	return &BoardHandler{boardService}
}

func (handler *BoardHandler) Routes() chi.Router {
	router := chi.NewRouter()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	router.Post("/upload/image", handler.uploadImage)

	return router
}

func (handler *BoardHandler) uploadImage(w http.ResponseWriter, r *http.Request) {
	userInfo, err := httputils.GetUserInfoData(r)
	if err != nil {
		httputils.ErrorJSON(w, err, http.StatusUnauthorized)
    return
	}

  media, err := httputils.GetUploadImageRequest(r)
  if err != nil {
    httputils.ErrorJSON(w, err, http.StatusBadRequest)
    return
  }

  result, err := handler.boardService.UploadImage(media, userInfo)
  if err != nil {
    httputils.ErrorJSON(w, err, http.StatusInternalServerError)
    return
  }

	httputils.WriteJSON(w, http.StatusOK, result)
}
