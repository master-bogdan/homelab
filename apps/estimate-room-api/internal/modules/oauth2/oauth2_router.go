package oauth2

import "github.com/gofiber/fiber/v2"

const (
	Authorize = "/authorize"
	Login     = "/login"
	Token     = "/token"
)

func SetupOauth2RoutesV1(router fiber.Router, controller Oauth2Controller) {
	auth := router.Group("/oauth2")

	auth.Get(Authorize, controller.Authorize)
	auth.Get(Login, controller.ShowLoginForm)
	auth.Post(Login, controller.Login)
	auth.Post(Token, controller.GetTokens)
}
