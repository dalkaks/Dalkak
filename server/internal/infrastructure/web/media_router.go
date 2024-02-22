package web

import (
	"dalkak/internal/core"
	mediadto "dalkak/pkg/dto/media"

	"github.com/gofiber/fiber/v3"
)

// prefix /media
func SetupMediaRoute(group fiber.Router, eventManager core.EventManager) {
	group.Post("/presigned", WarpHandler(func(c fiber.Ctx) interface{} {
		user, err := GetUserInfoFromContext(c, true)
		if err != nil {
			return err
		}

		req := new(mediadto.CreateTempMediaRequest)
		err = BindAndValidate(c, req)
		if err != nil {
			return err
		}

		return PublishAndWaitResponse(eventManager, "post.media.presigned", user, req)
	}))
}
