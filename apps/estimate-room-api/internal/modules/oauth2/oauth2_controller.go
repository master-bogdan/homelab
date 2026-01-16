package oauth2

// import (
// 	"fmt"
// 	"log/slog"
//
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/master-bogdan/clear-cash-api/internal/infra/db/postgresql/repositories"
// 	oauth2_dto "github.com/master-bogdan/clear-cash-api/internal/modules/oauth2/dto"
// 	oauth2_utils "github.com/master-bogdan/clear-cash-api/internal/modules/oauth2/utils"
// 	"github.com/master-bogdan/clear-cash-api/internal/pkg/utils"
// 	// "github.com/gofiber/fiber/v2"
// 	//"github.com/master-bogdan/clear-cash-api/internal/infra/db/postgresql/repositories"
// 	// oauth2_dto "github.com/master-bogdan/clear-cash-api/internal/modules/oauth2/dto"
// 	// oauth2_utils "github.com/master-bogdan/clear-cash-api/internal/modules/oauth2/utils"
// 	// "github.com/master-bogdan/clear-cash-api/internal/pkg/utils"
// )
//
// type Oauth2Controller interface {
// 	Authorize(ctx *fiber.Ctx) error
// 	Login(ctx *fiber.Ctx) error
// 	ShowLoginForm(ctx *fiber.Ctx) error
// 	GetTokens(ctx *fiber.Ctx) error
// }
//
// type oauth2Controller struct {
// 	service Oauth2Service
// 	logger  *slog.Logger
// }
//
// func NewOauth2Controller(
// 	oauth2Service Oauth2Service,
// 	logger *slog.Logger,
// ) Oauth2Controller {
// 	return &oauth2Controller{
// 		service: oauth2Service,
// 		logger:  logger,
// 	}
// }
//
// func (c *oauth2Controller) Authorize(ctx *fiber.Ctx) error {
// 	query := &oauth2_dto.AuthorizeQueryDTO{}
//
// 	err := ctx.QueryParser(query)
// 	if err != nil {
// 		c.logger.Error("invalid_query_params")
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "invalid_query_params"})
// 	}
//
// 	err = query.Validate()
// 	if err != nil {
// 		c.logger.Error(err.Error())
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	err = c.service.ValidateClient(query)
// 	if err != nil {
// 		c.logger.Error(err.Error())
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	XSessionIDHeader := ctx.Get("X-Session-Id")
//
// 	sessionID := ctx.Cookies("session_id")
//
// 	if XSessionIDHeader != "" {
// 		sessionID = XSessionIDHeader
// 	}
//
// 	userID := ""
//
// 	if sessionID != "" {
// 		userID, err = c.service.GetLoggedInUserID(sessionID)
// 	}
//
// 	if userID == "" || err != nil {
// 		c.logger.Warn("Session not found")
//
// 		loginRedirect := "/api/v1/oauth2/login?" + ctx.Context().QueryArgs().String()
// 		return ctx.Redirect(loginRedirect, fiber.StatusFound)
// 	}
//
// 	createAuthCodeDTO := &oauth2_dto.CreateOauthCodeDTO{
// 		ClientID:            query.ClientID,
// 		UserID:              userID,
// 		OidcSessionID:       sessionID,
// 		RedirectURI:         query.RedirectURI,
// 		CodeChallenge:       query.CodeChallenge,
// 		CodeChallengeMethod: query.CodeChallengeMethod,
// 		Scopes:              query.Scopes,
// 	}
// 	err = createAuthCodeDTO.Validate()
// 	if err != nil {
// 		c.logger.Error(err.Error())
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	authCode, err := c.service.CreateAuthCode(createAuthCodeDTO)
// 	if err != nil {
// 		c.logger.Error(err.Error())
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	redirectTo := query.RedirectURI + "?code=" + authCode
// 	if query.State != "" {
// 		redirectTo += "&state=" + query.State
// 	}
//
// 	c.logger.Info(fmt.Sprintf("Redirect to: %s", redirectTo))
//
// 	return ctx.Redirect(redirectTo, fiber.StatusFound)
// }
//
// func (c *oauth2Controller) ShowLoginForm(ctx *fiber.Ctx) error {
// 	query := &oauth2_dto.AuthorizeQueryDTO{}
//
// 	err := ctx.QueryParser(query)
// 	if err != nil {
// 		c.logger.Error("invalid query params")
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "invalid_query_params"})
// 	}
//
// 	params, err := utils.StructToMap(query)
// 	if err != nil {
// 		c.logger.Error("invalid query params")
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "invalid_query_params"})
// 	}
//
// 	ctx.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
// 	html := oauth2_utils.CreateLoginHtml(params)
//
// 	return ctx.SendString(html)
// }
//
// func (c *oauth2Controller) Login(ctx *fiber.Ctx) error {
// 	loginDTO := &oauth2_dto.LoginDTO{}
// 	err := ctx.BodyParser(loginDTO)
// 	if err != nil {
// 		c.logger.Error(fmt.Sprintf("Invalid body %s", err.Error()))
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	err = loginDTO.Validate()
// 	if err != nil {
// 		c.logger.Error(fmt.Sprintf("Invalid loginDTO %s", err.Error()))
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	userDTO := &oauth2_dto.UserDTO{
// 		Email:    loginDTO.Email,
// 		Password: loginDTO.Password,
// 	}
// 	err = userDTO.Validate()
// 	if err != nil {
// 		c.logger.Error(fmt.Sprintf("Invalid userDTO %s", err.Error()))
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	userID, err := c.service.AuthenticateUser(userDTO)
// 	if err != nil {
// 		// If user not found, register
// 		if err == repositories.ErrUserNotFound {
// 			c.logger.Warn("User not found")
//
// 			userID, err = c.service.RegisterUser(userDTO)
//
// 			if err != nil || userID == "" {
// 				c.logger.Error(err.Error())
// 				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 					"errors":  "Registration Failed",
// 					"details": err.Error(),
// 				})
// 			}
// 		}
// 	}
//
// 	oidcSessionDTO := &oauth2_dto.CreateOidcSessionDTO{
// 		UserID:   userID,
// 		ClientID: loginDTO.ClientID,
// 		Nonce:    loginDTO.Nonce,
// 	}
//
// 	err = oidcSessionDTO.Validate()
// 	if err != nil {
// 		c.logger.Error(fmt.Sprintf("Invalid oidcSessionDTO %s", err.Error()))
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	oidcSessionID, err := c.service.CreateOidcSession(oidcSessionDTO)
// 	if err != nil {
// 		c.logger.Error(err.Error())
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	ctx.Cookie(&fiber.Cookie{
// 		Name:     "session_id",
// 		Value:    oidcSessionID,
// 		HTTPOnly: true,
// 		Secure:   ctx.Protocol() == "https",
// 		Path:     "/",
// 	})
//
// 	// create auth code
// 	createAuthCodeDTO := &oauth2_dto.CreateOauthCodeDTO{
// 		ClientID:            loginDTO.ClientID,
// 		UserID:              userID,
// 		OidcSessionID:       oidcSessionID,
// 		RedirectURI:         loginDTO.RedirectURI,
// 		CodeChallenge:       loginDTO.CodeChallenge,
// 		CodeChallengeMethod: loginDTO.CodeChallengeMethod,
// 		Scopes:              loginDTO.Scopes,
// 	}
// 	err = createAuthCodeDTO.Validate()
// 	if err != nil {
// 		c.logger.Error(fmt.Sprintf("Invalid createAuthCodeDTO %s", err.Error()))
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	authCode, err := c.service.CreateAuthCode(createAuthCodeDTO)
// 	if err != nil {
// 		c.logger.Error(err.Error())
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	// Redirect to client with code and state
// 	redirectTo := loginDTO.RedirectURI + "?code=" + authCode
// 	if loginDTO.State != "" {
// 		redirectTo += "&state=" + loginDTO.State
// 	}
//
// 	c.logger.Info(fmt.Sprintf("redirecting to %s", redirectTo))
//
// 	return ctx.Redirect(redirectTo, fiber.StatusFound)
// }
//
// func (c *oauth2Controller) GetTokens(ctx *fiber.Ctx) error {
// 	body := &oauth2_dto.GetTokenDTO{}
//
// 	err := ctx.BodyParser(body)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	err = body.Validate()
// 	if err != nil {
// 		c.logger.Error(fmt.Sprintf("Invalid query %s", err.Error()))
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"errors": err.Error(),
// 		})
// 	}
//
// 	switch body.GrantType {
// 	case "authorization_code":
// 		tokens, err := c.service.GetAuthorizationTokens(body)
// 		if err != nil {
// 			c.logger.Error(err.Error())
// 			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": err.Error()})
// 		}
//
// 		return ctx.Status(fiber.StatusOK).JSON(tokens)
// 	case "refresh_token":
// 		tokens, err := c.service.GetRefreshTokens(body)
// 		if err != nil {
// 			c.logger.Error(err.Error())
// 			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"errors": err.Error()})
// 		}
//
// 		return ctx.Status(fiber.StatusOK).JSON(tokens)
// 	default:
// 		c.logger.Error("Unsupported grant_type")
// 		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "Unsupported grant_type"})
// 	}
// }
