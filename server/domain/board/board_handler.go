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
	// 권한 체크
	userInfo, err := httputils.GetUserInfoData(r)
	if err != nil {
		httputils.ErrorJSON(w, err, http.StatusUnauthorized)
	}

  content, err := httputils.GetUploadImageRequest(r)
  if err != nil {
    httputils.ErrorJSON(w, err, http.StatusBadRequest)
  }

	// 이미지 업로드
  // handler.boardService.UploadImage(content, userInfo)

	// 이미지 업로드 결과 반환
	httputils.WriteJSON(w, http.StatusOK, userInfo)
}
