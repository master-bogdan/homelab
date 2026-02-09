package oauth2

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql/repositories"
	oauth2dto "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/dto"
	oauth2utils "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/utils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
)

type Oauth2Controller interface {
	Authorize(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	ShowLoginForm(w http.ResponseWriter, r *http.Request)
	GetTokens(w http.ResponseWriter, r *http.Request)
	GithubLogin(w http.ResponseWriter, r *http.Request)
	GithubCallback(w http.ResponseWriter, r *http.Request)
}

type oauth2Controller struct {
	service       Oauth2Service
	logger        *slog.Logger
	github        oauth2utils.GithubConfig
	stateTokenKey []byte
	httpClient    *http.Client
}

func NewOauth2Controller(
	oauth2Service Oauth2Service,
	github oauth2utils.GithubConfig,
) Oauth2Controller {
	return &oauth2Controller{
		service:       oauth2Service,
		logger:        logger.L().With(slog.String("module", "oauth")),
		github:        github,
		stateTokenKey: []byte(github.StateSecret),
		httpClient:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *oauth2Controller) Authorize(w http.ResponseWriter, r *http.Request) {
	query := parseAuthorizeQuery(r)

	err := query.Validate()
	if err != nil {
		c.logger.Error(err.Error())
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.ValidateClient(query)
	if err != nil {
		c.logger.Error(err.Error())
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	sessionID := ""
	if cookie, err := r.Cookie("session_id"); err == nil {
		sessionID = cookie.Value
	}

	xSessionIDHeader := r.Header.Get("X-Session-Id")
	if xSessionIDHeader != "" {
		sessionID = xSessionIDHeader
	}

	userID := ""
	if sessionID != "" {
		userID, err = c.service.GetLoggedInUserID(sessionID)
	}

	if userID == "" || err != nil {
		c.logger.Warn("session not found")
		loginRedirect := "/api/v1/oauth2/login?" + r.URL.RawQuery
		http.Redirect(w, r, loginRedirect, http.StatusFound)
		return
	}

	createAuthCodeDTO := &oauth2dto.CreateOauthCodeDTO{
		ClientID:            query.ClientID,
		UserID:              userID,
		OidcSessionID:       sessionID,
		RedirectURI:         query.RedirectURI,
		CodeChallenge:       query.CodeChallenge,
		CodeChallengeMethod: query.CodeChallengeMethod,
		Scopes:              query.Scopes,
	}
	if err = createAuthCodeDTO.Validate(); err != nil {
		c.logger.Error(err.Error())
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	authCode, err := c.service.CreateAuthCode(createAuthCodeDTO)
	if err != nil {
		c.logger.Error(err.Error())
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	redirectTo := query.RedirectURI + "?code=" + authCode
	if query.State != "" {
		redirectTo += "&state=" + query.State
	}

	c.logger.Info(fmt.Sprintf("redirect to: %s", redirectTo))

	http.Redirect(w, r, redirectTo, http.StatusFound)
}

func (c *oauth2Controller) ShowLoginForm(w http.ResponseWriter, r *http.Request) {
	query := parseAuthorizeQuery(r)

	params, err := utils.StructToMap(query)
	if err != nil {
		c.logger.Error("invalid query params")
		utils.WriteResponseError(w, http.StatusBadRequest, "invalid_query_params")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := oauth2utils.CreateLoginHtml(params)
	_, _ = w.Write([]byte(html))
}

func (c *oauth2Controller) Login(w http.ResponseWriter, r *http.Request) {
	loginDTO, err := parseLoginForm(r)
	if err != nil {
		c.logger.Error(fmt.Sprintf("invalid body %s", err.Error()))
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = loginDTO.Validate(); err != nil {
		c.logger.Error(fmt.Sprintf("invalid loginDTO %s", err.Error()))
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	validateClientDTO := &oauth2dto.AuthorizeQueryDTO{
		ClientID:            loginDTO.ClientID,
		RedirectURI:         loginDTO.RedirectURI,
		ResponseType:        loginDTO.ResponseType,
		Scopes:              loginDTO.Scopes,
		State:               loginDTO.State,
		CodeChallenge:       loginDTO.CodeChallenge,
		CodeChallengeMethod: loginDTO.CodeChallengeMethod,
		Nonce:               loginDTO.Nonce,
	}
	if err = validateClientDTO.Validate(); err != nil {
		c.logger.Error(fmt.Sprintf("invalid client query %s", err.Error()))
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err = c.service.ValidateClient(validateClientDTO); err != nil {
		c.logger.Error(err.Error())
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	userDTO := &oauth2dto.UserDTO{
		Email:    loginDTO.Email,
		Password: loginDTO.Password,
	}
	if err = userDTO.Validate(); err != nil {
		c.logger.Error(fmt.Sprintf("invalid userDTO %s", err.Error()))
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := c.service.AuthenticateUser(userDTO)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			c.logger.Warn("user not found")
			userID, err = c.service.RegisterUser(userDTO)
			if err != nil || userID == "" {
				c.logger.Error(err.Error())
				utils.WriteResponseError(w, http.StatusInternalServerError, "registration failed")
				return
			}
		} else if err == ErrInvalidCredentials {
			c.logger.Warn("invalid credentials")
			utils.WriteResponseError(w, http.StatusUnauthorized, "invalid credentials")
			return
		} else {
			c.logger.Error(err.Error())
			utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	oidcSessionDTO := &oauth2dto.CreateOidcSessionDTO{
		UserID:   userID,
		ClientID: loginDTO.ClientID,
		Nonce:    loginDTO.Nonce,
	}

	if err = oidcSessionDTO.Validate(); err != nil {
		c.logger.Error(fmt.Sprintf("invalid oidcSessionDTO %s", err.Error()))
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	oidcSessionID, err := c.service.CreateOidcSession(oidcSessionDTO)
	if err != nil {
		c.logger.Error(err.Error())
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    oidcSessionID,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		Path:     "/",
	})

	createAuthCodeDTO := &oauth2dto.CreateOauthCodeDTO{
		ClientID:            loginDTO.ClientID,
		UserID:              userID,
		OidcSessionID:       oidcSessionID,
		RedirectURI:         loginDTO.RedirectURI,
		CodeChallenge:       loginDTO.CodeChallenge,
		CodeChallengeMethod: loginDTO.CodeChallengeMethod,
		Scopes:              loginDTO.Scopes,
	}
	if err = createAuthCodeDTO.Validate(); err != nil {
		c.logger.Error(fmt.Sprintf("invalid createAuthCodeDTO %s", err.Error()))
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	authCode, err := c.service.CreateAuthCode(createAuthCodeDTO)
	if err != nil {
		c.logger.Error(err.Error())
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	redirectTo := loginDTO.RedirectURI + "?code=" + authCode
	if loginDTO.State != "" {
		redirectTo += "&state=" + loginDTO.State
	}

	c.logger.Info(fmt.Sprintf("redirecting to %s", redirectTo))

	http.Redirect(w, r, redirectTo, http.StatusFound)
}

func (c *oauth2Controller) GetTokens(w http.ResponseWriter, r *http.Request) {
	body, err := parseTokenForm(r)
	if err != nil {
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = body.Validate(); err != nil {
		c.logger.Error(fmt.Sprintf("invalid query %s", err.Error()))
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch body.GrantType {
	case "authorization_code":
		tokens, err := c.service.GetAuthorizationTokens(r.Context(), body)
		if err != nil {
			c.logger.Error(err.Error())
			utils.WriteResponseError(w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.WriteResponse(w, http.StatusOK, tokens)
		return
	case "refresh_token":
		tokens, err := c.service.GetRefreshTokens(r.Context(), body)
		if err != nil {
			c.logger.Error(err.Error())
			utils.WriteResponseError(w, http.StatusUnauthorized, err.Error())
			return
		}

		utils.WriteResponse(w, http.StatusOK, tokens)
		return
	default:
		c.logger.Error("unsupported grant_type")
		utils.WriteResponseError(w, http.StatusBadRequest, "unsupported grant_type")
	}
}

// GithubLogin godoc
// @Summary Login with GitHub
// @Description Redirects to GitHub OAuth authorize URL.
// @Tags oauth2
// @Router /api/v1/oauth2/github/login [get]
func (c *oauth2Controller) GithubLogin(w http.ResponseWriter, r *http.Request) {
	if !c.github.IsConfigured() {
		utils.WriteResponseError(w, http.StatusInternalServerError, "github oauth is not configured")
		return
	}

	query := parseAuthorizeQuery(r)
	if err := query.Validate(); err != nil {
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.service.ValidateClient(query); err != nil {
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	stateToken, err := oauth2utils.CreateGithubState(c.stateTokenKey, oauth2utils.AuthorizeQueryFromDTO(query))
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	redirectURL := oauth2utils.BuildGithubAuthorizeURL(c.github, stateToken)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// GithubCallback godoc
// @Summary GitHub OAuth callback
// @Description Handles GitHub OAuth callback and issues internal auth code.
// @Tags oauth2
// @Router /api/v1/oauth2/github/callback [get]
func (c *oauth2Controller) GithubCallback(w http.ResponseWriter, r *http.Request) {
	if !c.github.IsConfigured() {
		utils.WriteResponseError(w, http.StatusInternalServerError, "github oauth is not configured")
		return
	}

	if r.URL.Query().Get("error") != "" {
		utils.WriteResponseError(w, http.StatusUnauthorized, "github authentication failed")
		return
	}

	code := r.URL.Query().Get("code")
	stateToken := r.URL.Query().Get("state")
	if code == "" || stateToken == "" {
		utils.WriteResponseError(w, http.StatusBadRequest, "code and state are required")
		return
	}

	state, err := oauth2utils.ParseGithubState(c.stateTokenKey, stateToken)
	if err != nil {
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	query := oauth2utils.AuthorizeQueryFromState(state).ToDTO()
	if err := query.Validate(); err != nil {
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.service.ValidateClient(query); err != nil {
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	accessToken, err := oauth2utils.ExchangeGithubCode(ctx, c.httpClient, c.github, code)
	if err != nil {
		utils.WriteResponseError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := oauth2utils.FetchGithubUser(ctx, c.httpClient, accessToken)
	if err != nil {
		utils.WriteResponseError(w, http.StatusUnauthorized, err.Error())
		return
	}

	emails, err := oauth2utils.FetchGithubEmails(ctx, c.httpClient, accessToken)
	if err != nil {
		c.logger.Warn("failed to fetch github emails", "err", err)
	}

	profile := oauth2utils.BuildGithubProfile(user, emails)
	userID, err := c.service.GetOrCreateUserFromGithub(profile)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			utils.WriteResponseError(w, http.StatusNotFound, "user not found")
			return
		}
		utils.WriteResponseError(w, http.StatusInternalServerError, "failed to login with github")
		return
	}

	oidcSessionDTO := &oauth2dto.CreateOidcSessionDTO{
		UserID:   userID,
		ClientID: query.ClientID,
		Nonce:    query.Nonce,
	}
	if err := oidcSessionDTO.Validate(); err != nil {
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	oidcSessionID, err := c.service.CreateOidcSession(oidcSessionDTO)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    oidcSessionID,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		Path:     "/",
	})

	createAuthCodeDTO := &oauth2dto.CreateOauthCodeDTO{
		ClientID:            query.ClientID,
		UserID:              userID,
		OidcSessionID:       oidcSessionID,
		RedirectURI:         query.RedirectURI,
		CodeChallenge:       query.CodeChallenge,
		CodeChallengeMethod: query.CodeChallengeMethod,
		Scopes:              query.Scopes,
	}
	if err := createAuthCodeDTO.Validate(); err != nil {
		utils.WriteResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	authCode, err := c.service.CreateAuthCode(createAuthCodeDTO)
	if err != nil {
		utils.WriteResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	redirectTo := query.RedirectURI + "?code=" + authCode
	if query.State != "" {
		redirectTo += "&state=" + query.State
	}

	http.Redirect(w, r, redirectTo, http.StatusFound)
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

func parseLoginForm(r *http.Request) (*oauth2dto.LoginDTO, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	return &oauth2dto.LoginDTO{
		Email:               r.Form.Get("email"),
		Password:            r.Form.Get("password"),
		ClientID:            r.Form.Get("client_id"),
		RedirectURI:         r.Form.Get("redirect_uri"),
		ResponseType:        r.Form.Get("response_type"),
		Scopes:              r.Form.Get("scopes"),
		State:               r.Form.Get("state"),
		CodeChallenge:       r.Form.Get("code_challenge"),
		CodeChallengeMethod: r.Form.Get("code_challenge_method"),
		Nonce:               r.Form.Get("nonce"),
	}, nil
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
