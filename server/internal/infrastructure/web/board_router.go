package web

import (
	"dalkak/internal/core"
	boarddto "dalkak/pkg/dto/board"

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
}
