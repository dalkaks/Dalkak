package web

import (
	"dalkak/internal/core"
	boarddto "dalkak/pkg/dto/board"
	responseutil "dalkak/pkg/utils/response"

	"github.com/gofiber/fiber/v3"
)

// prefix /board
func SetupBoardRoute(group fiber.Router, eventManager core.EventManager) {
	group.Post("/", WarpHandler(func(c fiber.Ctx) interface{} {
		user, err := GetUserInfoFromContext(c, true)
		if err != nil {
			return err
		}

		req := new(boarddto.CreateBoardRequest)
		err = BindAndValidate(c, req)
		if err != nil {
			return err
		}

		return PublishAndWaitResponse(eventManager, "post.board", user, req)
	}))

	group.Get("/list/processing", WarpHandler(func(c fiber.Ctx) interface{} {
		user, err := GetUserInfoFromContext(c, true)
		if err != nil {
			return err
		}

		req := new(boarddto.GetBoardListProcessingRequest)
		err = BindAndValidate(c, req)
		if err != nil {
			return err
		}

		return PublishAndWaitResponse(eventManager, "get.board.list.processing", user, req)
	}))

	group.Delete("/:id", WarpHandler(func(c fiber.Ctx) interface{} {
		user, err := GetUserInfoFromContext(c, true)
		if err != nil {
			return err
		}

		req := new(boarddto.DeleteBoardRequest)
		req.Id = c.Params("id")
		if req.Id == "" {
			return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
		}

		return PublishAndWaitResponse(eventManager, "delete.board", user, req)
	}))
}
