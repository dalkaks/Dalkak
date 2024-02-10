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

	router.Post("/image/presigned", handler.createPresignedURL)

	return router
}

func (handler *BoardHandler) createPresignedURL(w http.ResponseWriter, r *http.Request) {
	userInfo, err := httputils.GetUserInfoData(r)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	var req payloads.BoardUploadMediaRequest
	err = httputils.ReadJSON(w, r, &req)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	result, err := handler.boardService.CreatePresignedURL(&req, userInfo)
	if err != nil {
		httputils.HandleAppError(w, err)
		return
	}

	httputils.WriteJSONAndHandleError(w, http.StatusOK, result, httputils.HandleAppError)
}
