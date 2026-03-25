package oauth2

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	oauth2dto "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/dto"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/httputils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type Oauth2Controller interface {
	Authorize(w http.ResponseWriter, r *http.Request)
	GetTokens(w http.ResponseWriter, r *http.Request)
}

type oauth2Controller struct {
	service         Oauth2Service
	authService     AuthService
	logger          *slog.Logger
	frontendBaseURL string
}

func NewOauth2Controller(
	oauth2Service Oauth2Service,
	authService AuthService,
	frontendBaseURL string,
) Oauth2Controller {
	return &oauth2Controller{
		service:         oauth2Service,
		authService:     authService,
		logger:          logger.L().With(slog.String("module", "oauth")),
		frontendBaseURL: frontendBaseURL,
	}
}

// Authorize godoc
// @Summary OAuth2 authorize
// @Description Validates the authorization request and either issues an auth code for an authenticated browser session or redirects the user to the frontend login page.
// @Tags oauth2
// @Produce plain
// @Param client_id query string true "OAuth client ID"
// @Param redirect_uri query string true "Registered redirect URI"
// @Param response_type query string true "Must be code"
// @Param scopes query string true "Space separated scopes"
// @Param state query string true "Opaque client state"
// @Param code_challenge query string true "PKCE challenge"
// @Param code_challenge_method query string true "PKCE method"
// @Param nonce query string true "OIDC nonce"
// @Success 302 {string} string "Redirect to frontend login or client redirect URI"
// @Failure 400 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/oauth2/authorize [get]
func (c *oauth2Controller) Authorize(w http.ResponseWriter, r *http.Request) {
	query := parseAuthorizeQuery(r)

	err := query.Validate()
	if err != nil {
		c.logger.Error(err.Error())
		httputils.WriteResponseError(w, apperrors.CreateHttpError(apperrors.ErrBadRequest, apperrors.HttpError{Detail: "invalid params"}))
		return
	}

	err = c.service.ValidateClient(query)
	if err != nil {
		c.logger.Error(err.Error())
		httputils.WriteResponseError(w, apperrors.CreateHttpError(apperrors.ErrBadRequest, apperrors.HttpError{Detail: err.Error()}))
		return
	}

	sessionID := ReadSessionID(r)
	if sessionID == "" {
		c.redirectToFrontendLogin(w, r)
		return
	}

	session, err := c.authService.GetOidcSessionByID(sessionID)
	if err != nil {
		c.logger.Warn("session not found", "err", err)
		c.redirectToFrontendLogin(w, r)
		return
	}

	if session.ClientID != query.ClientID || session.Nonce != query.Nonce {
		oidcSessionID, createErr := c.service.CreateOidcSession(&oauth2dto.CreateOidcSessionDTO{
			UserID:   session.UserID,
			ClientID: query.ClientID,
			Nonce:    query.Nonce,
		})
		if createErr != nil {
			c.logger.Error(createErr.Error())
			httputils.WriteResponseError(w, apperrors.CreateHttpError(apperrors.ErrInternal, apperrors.HttpError{Detail: createErr.Error()}))
			return
		}

		sessionID = oidcSessionID
		http.SetCookie(w, SessionCookie(oidcSessionID, r.TLS != nil))
	}

	createAuthCodeDTO := &oauth2dto.CreateOauthCodeDTO{
		ClientID:            query.ClientID,
		UserID:              session.UserID,
		OidcSessionID:       sessionID,
		RedirectURI:         query.RedirectURI,
		CodeChallenge:       query.CodeChallenge,
		CodeChallengeMethod: query.CodeChallengeMethod,
		Scopes:              query.Scopes,
	}
	if err = createAuthCodeDTO.Validate(); err != nil {
		c.logger.Error(err.Error())
		httputils.WriteResponseError(w, apperrors.CreateHttpError(apperrors.ErrBadRequest, apperrors.HttpError{Detail: err.Error()}))
		return
	}

	authCode, err := c.service.CreateAuthCode(createAuthCodeDTO)
	if err != nil {
		c.logger.Error(err.Error())
		httputils.WriteResponseError(w, apperrors.CreateHttpError(apperrors.ErrInternal, apperrors.HttpError{Detail: err.Error()}))
		return
	}

	redirectTo := query.RedirectURI + "?code=" + authCode
	if query.State != "" {
		redirectTo += "&state=" + query.State
	}

	c.logger.Info(fmt.Sprintf("redirect to: %s", redirectTo))

	http.Redirect(w, r, redirectTo, http.StatusFound)
}

// GetTokens godoc
// @Summary Exchange OAuth2 token
// @Description Exchanges an authorization code or refresh token for access, refresh, and ID tokens.
// @Tags oauth2
// @Accept x-www-form-urlencoded
// @Produce json
// @Param grant_type formData string true "authorization_code or refresh_token"
// @Param code formData string false "Authorization code"
// @Param code_verifier formData string false "PKCE verifier"
// @Param client_id formData string true "OAuth client ID"
// @Param redirect_uri formData string false "Redirect URI"
// @Param refresh_token formData string false "Refresh token"
// @Success 200 {object} oauth2dto.TokenResponseDTO
// @Failure 400 {object} apperrors.HttpError
// @Failure 401 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/oauth2/token [post]
func (c *oauth2Controller) GetTokens(w http.ResponseWriter, r *http.Request) {
	body, err := parseTokenForm(r)
	if err != nil {
		httputils.WriteResponseError(w, apperrors.CreateHttpError(apperrors.ErrBadRequest, apperrors.HttpError{Detail: err.Error()}))
		return
	}

	if err = body.Validate(); err != nil {
		c.logger.Error(fmt.Sprintf("invalid query %s", err.Error()))
		httputils.WriteResponseError(w, apperrors.CreateHttpError(apperrors.ErrBadRequest, apperrors.HttpError{Detail: err.Error()}))
		return
	}

	switch body.GrantType {
	case "authorization_code":
		tokens, err := c.service.GetAuthorizationTokens(r.Context(), body)
		if err != nil {
			c.logger.Error(err.Error())
			httputils.WriteResponseError(w, apperrors.CreateHttpError(apperrors.ErrUnauthorized, apperrors.HttpError{Detail: err.Error()}))
			return
		}

		httputils.WriteResponse(w, tokens, httputils.WriteResponseOptions{Status: http.StatusOK})
		return
	case "refresh_token":
		tokens, err := c.service.GetRefreshTokens(r.Context(), body)
		if err != nil {
			c.logger.Error(err.Error())
			httputils.WriteResponseError(w, apperrors.CreateHttpError(apperrors.ErrUnauthorized, apperrors.HttpError{Detail: err.Error()}))
			return
		}

		httputils.WriteResponse(w, tokens, httputils.WriteResponseOptions{Status: http.StatusOK})
		return
	default:
		c.logger.Error("unsupported grant_type")
		httputils.WriteResponseError(w, apperrors.CreateHttpError(apperrors.ErrBadRequest, apperrors.HttpError{Detail: "unsupported grant_type"}))
	}
}

func parseAuthorizeQuery(r *http.Request) *oauth2dto.AuthorizeQueryDTO {
	q := r.URL.Query()
	return &oauth2dto.AuthorizeQueryDTO{
		ClientID:            q.Get("client_id"),
		RedirectURI:         q.Get("redirect_uri"),
		ResponseType:        q.Get("response_type"),
		Scopes:              q.Get("scopes"),
		State:               q.Get("state"),
		CodeChallenge:       q.Get("code_challenge"),
		CodeChallengeMethod: q.Get("code_challenge_method"),
		Nonce:               q.Get("nonce"),
	}
}

func parseTokenForm(r *http.Request) (*oauth2dto.GetTokenDTO, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	return &oauth2dto.GetTokenDTO{
		GrantType:    r.Form.Get("grant_type"),
		CodeVerifier: r.Form.Get("code_verifier"),
		Code:         r.Form.Get("code"),
		ClientID:     r.Form.Get("client_id"),
		RedirectURI:  r.Form.Get("redirect_uri"),
		RefreshToken: r.Form.Get("refresh_token"),
	}, nil
}

func (c *oauth2Controller) redirectToFrontendLogin(w http.ResponseWriter, r *http.Request) {
	loginURL, err := buildFrontendLoginURL(c.frontendBaseURL, currentRequestURL(r))
	if err != nil {
		c.logger.Error("failed to build frontend login redirect", "err", err)
		httputils.WriteResponseError(w, apperrors.CreateHttpError(
			apperrors.ErrInternal,
			apperrors.HttpError{Detail: "frontend auth redirect is not configured"},
		))
		return
	}

	http.Redirect(w, r, loginURL, http.StatusFound)
}

func buildFrontendLoginURL(frontendBaseURL, continueURL string) (string, error) {
	trimmed := strings.TrimRight(strings.TrimSpace(frontendBaseURL), "/")
	if trimmed == "" {
		return "", fmt.Errorf("frontend base url is required")
	}

	loginURL, err := url.Parse(trimmed + "/login")
	if err != nil {
		return "", err
	}

	query := loginURL.Query()
	query.Set("continue", continueURL)
	loginURL.RawQuery = query.Encode()

	return loginURL.String(), nil
}

func currentRequestURL(r *http.Request) string {
	if r == nil || r.URL == nil {
		return ""
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	if forwardedProto := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-Proto"), ",")[0]); forwardedProto != "" {
		scheme = forwardedProto
	}

	copyURL := *r.URL
	if !copyURL.IsAbs() {
		copyURL.Scheme = scheme
		copyURL.Host = r.Host
	}

	return copyURL.String()
}
