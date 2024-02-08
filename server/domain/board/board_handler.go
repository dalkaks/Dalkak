package board

import (
	"dalkak/pkg/interfaces"
	"dalkak/pkg/payloads"
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

	router.Post("/image/presigned", handler.createPresignedURL)

	return router
}

func (handler *BoardHandler) createPresignedURL(w http.ResponseWriter, r *http.Request) {
	userInfo, err := httputils.GetUserInfoData(r)
	if err != nil {
		httputils.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	req, err := httputils.GetRequestData[payloads.BoardUploadMediaRequest](r)
	if err != nil {
		httputils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	result, err := handler.boardService.CreatePresignedURL(req, userInfo)
	if err != nil {
		// Todo: error handle
		httputils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err := httputils.WriteJSON(w, http.StatusOK, result); err != nil {
		httputils.ErrorJSON(w, err, http.StatusInternalServerError)
	}
}
