package web

import (
	"dalkak/internal/core"
	"dalkak/internal/infrastructure/eventbus"
	appdto "dalkak/pkg/dto/app"
	responseutil "dalkak/pkg/utils/response"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/idempotency"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

const DefaultTimeout = 20 * time.Second

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type HandleFunc func(c fiber.Ctx) interface{}

func NewRouter(mode string, origin string, infra *core.Infra) *fiber.App {
	router := fiber.New()

	router.Use(recover.New())
	router.Use(logger.New())

	router.Use(helmet.New())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     origin,
		AllowMethods:     "OPTIONS,GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept, X-CSRF-Token, Authorization, x-client-id",
		AllowCredentials: true,
	}))
	if mode == "PROD" {
		router.Use(csrf.New())
	}

	router.Use(idempotency.New())
	router.Use(etag.New(etag.Config{
		Weak: true,
	}))

	router.Use(func(c fiber.Ctx) error {
		return getAuthUserMiddleware(c, infra.Keymanager)
	})

	router.Get("/", func(c fiber.Ctx) error {
		return c.SendString("ok")
	})

	router.Get("/test", WarpHandler(func(c fiber.Ctx) interface{} {
		return responseutil.NewAppData(origin)
	}))

	userRoute := router.Group("/user")
	SetupUserRoute(userRoute, infra.Keymanager, infra.EventManager)

	mediaRoute := router.Group("/media")
	SetupMediaRoute(mediaRoute, infra.EventManager)

	// Default not found handler
	router.All("*", WarpHandler(func(c fiber.Ctx) interface{} {
		return responseutil.NewAppError(responseutil.ErrCodeNotFound, responseutil.ErrMsgRequestNotFound)
	}))

	return router
}

func WarpHandler(handler HandleFunc) fiber.Handler {
	return func(c fiber.Ctx) error {
		resp := handler(c)
		responseutil.WriteToResponse(c, resp)
		return nil
	}
}

func GetUserInfoFromContext(c fiber.Ctx, requireUserInfo bool) (*appdto.UserInfo, error) {
	userInfo, ok := c.Locals("user").(appdto.UserInfo)
	if !ok {
		if requireUserInfo {
			return nil, responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgRequestUnauth)
		}
		return nil, nil
	}
	return &userInfo, nil
}

func BindAndValidate(c fiber.Ctx, req interface{}) error {
	c.Bind().Body(req)
	c.Bind().Query(req)

	err := validate.Struct(req)
	if err != nil {
		return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid, err)
	}
	return nil
}

func PublishAndWaitResponse(eventManager core.EventManager, eventType string, userInfo *appdto.UserInfo, payload interface{}, timeout ...time.Duration) interface{} {
	actualTimeout := DefaultTimeout
	if len(timeout) > 0 {
		actualTimeout = timeout[0]
	}

	responseChan := make(chan appdto.Response, 1)
	defer close(responseChan)

	eventManager.Publish(eventbus.Event{
		Type:         eventType,
		UserInfo:     userInfo,
		ResponseChan: responseChan,
		Payload:      payload,
	})
	return waitForResponse(responseChan, actualTimeout)
}

func waitForResponse(responseChan chan appdto.Response, timeout time.Duration) interface{} {
	select {
	case resp := <-responseChan:
		if resp.Error != nil {
			return responseutil.NewAppError(responseutil.ErrCodeInternal, resp.Error.Error())
		} else if appData, ok := resp.Data.(*responseutil.AppData); ok {
			return appData
		}
		return responseutil.NewAppData(resp.Data, responseutil.DataCodeSuccess)
	case <-time.After(timeout):
		return responseutil.NewAppError(responseutil.ErrCodeTimeout, responseutil.ErrMsgRequestTimeout)
	}
}

func getAuthUserMiddleware(c fiber.Ctx, keyManager core.KeyManager) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Next()
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		err := responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgTokenParseFailed)
		responseutil.WriteToResponse(c, err)
		return err
	}

	token := headerParts[1]
	sub, err := keyManager.ParseTokenWithPublicKey(token)
	if err != nil {
		responseutil.WriteToResponse(c, err)
		return err
	}

	c.Locals("user", appdto.UserInfo{WalletAddress: sub})
	return c.Next()
}
