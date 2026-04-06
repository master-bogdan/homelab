package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	authdto "github.com/master-bogdan/estimate-room-api/internal/modules/auth/dto"
	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
	apperrors "github.com/master-bogdan/estimate-room-api/internal/pkg/apperrors"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/httputils"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

type AuthController interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	ForgotPassword(w http.ResponseWriter, r *http.Request)
	ValidateResetPasswordToken(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	GetSession(w http.ResponseWriter, r *http.Request)
	GithubLogin(w http.ResponseWriter, r *http.Request)
	GithubCallback(w http.ResponseWriter, r *http.Request)
}

type authController struct {
	service           AuthService
	logger            *slog.Logger
	trustProxyHeaders bool
}

func NewAuthController(service AuthService, trustProxyHeaders bool) AuthController {
	return &authController{
		service:           service,
		logger:            logger.L().With(slog.String("controller", "auth")),
		trustProxyHeaders: trustProxyHeaders,
	}
}

// Login godoc
// @Summary Login with email and password
// @Description Authenticates a local user, creates a browser session, and returns the authenticated user payload.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body authdto.LoginDTO true "Credentials and OAuth continue URL"
// @Success 200 {object} authdto.SessionResponse
// @Failure 400 {object} apperrors.HttpError
// @Failure 401 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/auth/login [post]
func (c *authController) Login(w http.ResponseWriter, r *http.Request) {
	dto := authdto.LoginDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}
	if err := dto.Validate(); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	logArgs := []any{
		"path", r.URL.Path,
		"email", maskEmailForLog(dto.Email),
	}
	logArgs = append(logArgs, safeContinueLogFields(dto.ContinueURL)...)

	logger.FromRequest(r, c.logger).Info("login dto accepted", logArgs...)

	user, sessionID, err := c.service.Login(r.Context(), &dto)
	if err != nil {
		c.writeAuthError(w, r, err)
		return
	}

	http.SetCookie(w, oauth2.Oauth2SessionCookie(sessionID, r, c.trustProxyHeaders))
	httputils.WriteResponse(w, formatSessionResponse(user))
}

// Register godoc
// @Summary Register a local account
// @Description Creates a local user, auto-signs them in, and returns the authenticated user payload.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body authdto.RegisterDTO true "Registration payload and OAuth continue URL"
// @Success 201 {object} authdto.SessionResponse
// @Failure 400 {object} apperrors.HttpError
// @Failure 409 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/auth/register [post]
func (c *authController) Register(w http.ResponseWriter, r *http.Request) {
	dto := authdto.RegisterDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}
	if err := dto.Validate(); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	logArgs := []any{
		"path", r.URL.Path,
		"email", maskEmailForLog(dto.Email),
		"display_name_provided", dto.DisplayName != "",
		"organization_provided", dto.Organization != "",
		"occupation_provided", dto.Occupation != "",
	}
	logArgs = append(logArgs, safeContinueLogFields(dto.ContinueURL)...)

	logger.FromRequest(r, c.logger).Info("register dto accepted", logArgs...)

	user, sessionID, err := c.service.Register(r.Context(), &dto)
	if err != nil {
		c.writeAuthError(w, r, err)
		return
	}

	http.SetCookie(w, oauth2.Oauth2SessionCookie(sessionID, r, c.trustProxyHeaders))
	httputils.WriteResponse(w, formatSessionResponse(user), httputils.WriteResponseOptions{Status: http.StatusCreated})
}

// ForgotPassword godoc
// @Summary Start password reset
// @Description Accepts an email and stores a one-time reset token when the account is eligible. The response is always generic.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body authdto.ForgotPasswordDTO true "Email address"
// @Success 200 {object} authdto.ForgotPasswordResponse
// @Failure 400 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/auth/forgot-password [post]
func (c *authController) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	dto := authdto.ForgotPasswordDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}
	if err := dto.Validate(); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	logger.FromRequest(r, c.logger).Info("forgot password dto accepted",
		"path", r.URL.Path,
		"email", maskEmailForLog(dto.Email),
	)

	if err := c.service.ForgotPassword(r.Context(), &dto); err != nil {
		c.writeError(w, r, apperrors.ErrInternal, "", err)
		return
	}

	httputils.WriteResponse(w, authdto.ForgotPasswordResponse{Submitted: true})
}

// ValidateResetPasswordToken godoc
// @Summary Validate password reset token
// @Description Returns whether a password reset token is valid, expired, invalid, or already used.
// @Tags auth
// @Produce json
// @Param token query string true "Reset token"
// @Success 200 {object} authdto.ResetPasswordValidationResponse
// @Failure 400 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/auth/reset-password/validate [get]
func (c *authController) ValidateResetPasswordToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		c.writeError(w, r, apperrors.ErrBadRequest, "token is required", nil)
		return
	}

	logger.FromRequest(r, c.logger).Info("reset password token validation requested",
		"path", r.URL.Path,
		"token_present", true,
		"token_length", len(token),
	)

	valid, reason, err := c.service.ValidateResetPasswordToken(r.Context(), token)
	if err != nil {
		c.writeError(w, r, apperrors.ErrInternal, "", err)
		return
	}

	httputils.WriteResponse(w, authdto.ResetPasswordValidationResponse{
		Valid:  valid,
		Reason: reason,
	})
}

// ResetPassword godoc
// @Summary Reset local account password
// @Description Resets the password for a valid token and revokes all active sessions and tokens for that user.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body authdto.ResetPasswordDTO true "Reset token and new password"
// @Success 200 {object} authdto.ResetPasswordResponse
// @Failure 400 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/auth/reset-password [post]
func (c *authController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	dto := authdto.ResetPasswordDTO{}
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}
	if err := dto.Validate(); err != nil {
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
		return
	}

	logger.FromRequest(r, c.logger).Info("reset password dto accepted",
		"path", r.URL.Path,
		"token_present", dto.Token != "",
		"token_length", len(dto.Token),
		"password_length", len(dto.Password),
	)

	if err := c.service.ResetPassword(r.Context(), &dto); err != nil {
		c.writeAuthError(w, r, err)
		return
	}

	http.SetCookie(w, oauth2.ExpiredOauth2SessionCookie(r, c.trustProxyHeaders))
	http.SetCookie(w, oauth2.ExpiredOauth2AccessTokenCookie(r, c.trustProxyHeaders))
	http.SetCookie(w, oauth2.ExpiredOauth2RefreshTokenCookie(r, c.trustProxyHeaders))

	httputils.WriteResponse(w, authdto.ResetPasswordResponse{Reset: true})
}

// Logout godoc
// @Summary Logout current session
// @Description Revokes the current browser session and clears auth cookies.
// @Tags auth
// @Produce json
// @Success 200 {object} authdto.LogoutResponse
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/auth/logout [post]
func (c *authController) Logout(w http.ResponseWriter, r *http.Request) {
	sessionID := oauth2.ReadOauth2SessionID(r)
	if err := c.service.Logout(r.Context(), sessionID); err != nil {
		c.writeError(w, r, apperrors.ErrInternal, "", err)
		return
	}

	http.SetCookie(w, oauth2.ExpiredOauth2SessionCookie(r, c.trustProxyHeaders))
	http.SetCookie(w, oauth2.ExpiredOauth2AccessTokenCookie(r, c.trustProxyHeaders))
	http.SetCookie(w, oauth2.ExpiredOauth2RefreshTokenCookie(r, c.trustProxyHeaders))

	httputils.WriteResponse(w, authdto.LogoutResponse{LoggedOut: true})
}

// GetSession godoc
// @Summary Current browser auth session
// @Description Returns the currently authenticated user for the active browser session, or authenticated=false when there is no active session.
// @Tags auth
// @Produce json
// @Success 200 {object} authdto.SessionResponse
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/auth/session [get]
func (c *authController) GetSession(w http.ResponseWriter, r *http.Request) {
	user, authenticated, err := c.service.GetSession(r)
	if err != nil {
		c.writeError(w, r, apperrors.ErrInternal, "", err)
		return
	}
	if !authenticated {
		httputils.WriteResponse(w, authdto.SessionResponse{
			Authenticated: false,
			User:          nil,
		})
		return
	}

	httputils.WriteResponse(w, formatSessionResponse(user))
}

// GithubLogin godoc
// @Summary Start GitHub login
// @Description Redirects the browser to GitHub OAuth using the provided OAuth continue URL.
// @Tags auth
// @Produce plain
// @Param continue query string true "OAuth continue URL"
// @Success 302 {string} string "Redirect to GitHub"
// @Failure 400 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/auth/github/login [get]
func (c *authController) GithubLogin(w http.ResponseWriter, r *http.Request) {
	continueURL := r.URL.Query().Get("continue")
	logArgs := []any{"path", r.URL.Path}
	logArgs = append(logArgs, safeContinueLogFields(continueURL)...)

	logger.FromRequest(r, c.logger).Info("github login requested", logArgs...)
	redirectURL, err := c.service.StartGithubLogin(continueURL)
	if err != nil {
		c.writeAuthError(w, r, err)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// GithubCallback godoc
// @Summary GitHub auth callback
// @Description Completes GitHub login, creates a browser session, and redirects back to the OAuth continue URL.
// @Tags auth
// @Produce plain
// @Success 302 {string} string "Redirect to continue URL"
// @Failure 400 {object} apperrors.HttpError
// @Failure 401 {object} apperrors.HttpError
// @Failure 500 {object} apperrors.HttpError
// @Router /api/v1/auth/github/callback [get]
func (c *authController) GithubCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("error") != "" {
		c.writeAuthError(w, r, ErrGithubAuthenticationFailed)
		return
	}

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		c.writeError(w, r, apperrors.ErrBadRequest, "code and state are required", nil)
		return
	}

	logger.FromRequest(r, c.logger).Info("github callback params accepted",
		"path", r.URL.Path,
		"code_present", true,
		"state_present", true,
	)

	continueURL, sessionID, err := c.service.HandleGithubCallback(r.Context(), code, state)
	if err != nil {
		c.writeAuthError(w, r, err)
		return
	}

	http.SetCookie(w, oauth2.Oauth2SessionCookie(sessionID, r, c.trustProxyHeaders))
	http.Redirect(w, r, continueURL, http.StatusFound)
}

func safeContinueLogFields(continueURL string) []any {
	trimmedContinueURL := strings.TrimSpace(continueURL)
	fields := []any{"continue_present", trimmedContinueURL != ""}
	if trimmedContinueURL == "" {
		return fields
	}

	parsedURL, err := url.Parse(trimmedContinueURL)
	if err == nil && parsedURL.Path != "" {
		fields = append(fields, "continue_path", parsedURL.Path)
	}

	return fields
}

func (c *authController) writeAuthError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrGithubAuthenticationFailed):
		c.writeError(w, r, apperrors.ErrUnauthorized, err.Error(), err)
	case errors.Is(err, ErrEmailAlreadyInUse):
		c.writeError(w, r, apperrors.ErrConflict, err.Error(), err)
	case errors.Is(err, ErrInvalidContinueURL),
		errors.Is(err, ErrInvalidResetToken),
		errors.Is(err, ErrExpiredResetToken),
		errors.Is(err, ErrUsedResetToken):
		c.writeError(w, r, apperrors.ErrBadRequest, err.Error(), err)
	case errors.Is(err, ErrGithubAuthNotConfigured):
		c.writeError(w, r, apperrors.ErrInternal, err.Error(), err)
	default:
		c.writeError(w, r, apperrors.ErrInternal, "", err)
	}
}

func (c *authController) writeError(w http.ResponseWriter, r *http.Request, errType error, detail string, cause error) {
	logArgs := []any{
		"path", r.URL.Path,
		"type", errType.Error(),
	}
	if detail != "" {
		logArgs = append(logArgs, "detail", detail)
	}
	if cause != nil {
		logArgs = append(logArgs, "err", cause)
	}

	logger.FromRequest(r, c.logger).Error("request failed", logArgs...)

	httputils.WriteResponseError(w, apperrors.CreateHttpError(
		errType,
		apperrors.HttpError{
			Detail:   detail,
			Instance: r.URL.Path,
		},
	))
}

func maskEmailForLog(email string) string {
	parts := strings.Split(strings.TrimSpace(email), "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "[REDACTED_EMAIL]"
	}

	switch len(parts[0]) {
	case 1:
		return "*@" + parts[1]
	case 2:
		return parts[0][:1] + "*@" + parts[1]
	default:
		return parts[0][:1] + "***@" + parts[1]
	}
}
