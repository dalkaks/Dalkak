package web

import (
	"dalkak/internal/core"
	userdto "dalkak/pkg/dto/user"
	responseutil "dalkak/pkg/utils/response"

	"github.com/gofiber/fiber/v3"
)

// prefix /user
func SetupUserRoute(group fiber.Router, keyManager core.KeyManager, eventManager core.EventManager) {
	group.Post("/auth", WarpHandler(func(c fiber.Ctx) interface{} {
		req := new(userdto.AuthAndSignUpRequest)
		err := BindAndValidate(c, req)
		if err != nil {
			return err
		}

		return PublishAndWaitResponse(eventManager, "post.user.auth", nil, req)
	}))

	group.Post("/refresh", WarpHandler(func(c fiber.Ctx) interface{} {
		refreshToken := c.Cookies("refreshToken")
		if refreshToken == "" {
			return responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgTokenParseFailed)
		}
		sub, err := keyManager.ParseTokenWithPublicKey(refreshToken)
		if err != nil {
			return err
		}

		req := &userdto.ReissueAccessTokenRequest{
			WalletAddress: sub,
		}

		return PublishAndWaitResponse(eventManager, "post.user.refresh", nil, req)
	}))

	group.Post("/media/presigned", WarpHandler(func(c fiber.Ctx) interface{} {
		user, err := GetUserInfoFromContext(c, true)
		if err != nil {
			return err
		}

		req := new(userdto.CreateTempMediaRequest)
		err = BindAndValidate(c, req)
		if err != nil {
			return err
		}

		return PublishAndWaitResponse(eventManager, "post.user.media.presigned", user, req)
	}))
}
