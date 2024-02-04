package board

import (
	"dalkak/pkg/interfaces"
	"net/http"
	"dalkak/pkg/utils/httputils"

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
  // 권한 체크
  userInfo, err := httputils.GetUserInfoData(r)
  if err != nil {
    httputils.ErrorJSON(w, err, http.StatusUnauthorized)
  }

  // 이미지 업로드

  // 이미지 업로드 결과 반환
}