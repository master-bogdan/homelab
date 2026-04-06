package oauth2

import (
	"net/http"
	"strings"
)

const (
	SessionCookieName      = "session_id"
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

func ReadSessionID(r *http.Request) string {
	if r == nil {
		return ""
	}

	if cookie, err := r.Cookie(SessionCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	if header := r.Header.Get("X-Session-Id"); header != "" {
		return header
	}

	return ""
}

func SessionCookie(sessionID string, r *http.Request, trustProxyHeaders bool) *http.Cookie {
	return cookieWithSettings(SessionCookieName, sessionID, 0, isSecureRequest(r, trustProxyHeaders))
}

func AccessTokenCookie(token string, r *http.Request, trustProxyHeaders bool) *http.Cookie {
	return cookieWithSettings(
		AccessTokenCookieName,
		token,
		int(accessTokenTTL.Seconds()),
		isSecureRequest(r, trustProxyHeaders),
	)
}

func RefreshTokenCookie(token string, r *http.Request, trustProxyHeaders bool) *http.Cookie {
	return cookieWithSettings(
		RefreshTokenCookieName,
		token,
		int(refreshTokenTTL.Seconds()),
		isSecureRequest(r, trustProxyHeaders),
	)
}

func ExpiredSessionCookie(r *http.Request, trustProxyHeaders bool) *http.Cookie {
	return expiredCookie(SessionCookieName, isSecureRequest(r, trustProxyHeaders))
}

func ExpiredAccessTokenCookie(r *http.Request, trustProxyHeaders bool) *http.Cookie {
	return expiredCookie(AccessTokenCookieName, isSecureRequest(r, trustProxyHeaders))
}

func ExpiredRefreshTokenCookie(r *http.Request, trustProxyHeaders bool) *http.Cookie {
	return expiredCookie(RefreshTokenCookieName, isSecureRequest(r, trustProxyHeaders))
}

func cookieWithSettings(name, value string, maxAge int, secure bool) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   secure,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	if maxAge > 0 {
		cookie.MaxAge = maxAge
	}

	return cookie
}

func expiredCookie(name string, secure bool) *http.Cookie {
	cookie := cookieWithSettings(name, "", 0, secure)
	cookie.MaxAge = -1

	return cookie
}

func isSecureRequest(r *http.Request, trustProxyHeaders bool) bool {
	if r == nil {
		return false
	}
	if r.TLS != nil {
		return true
	}
	if !trustProxyHeaders {
		return false
	}

	return strings.EqualFold(firstForwardedValue(r.Header.Get("X-Forwarded-Proto")), "https")
}

func firstForwardedValue(value string) string {
	head := strings.TrimSpace(strings.Split(value, ",")[0])
	if head == "" {
		return ""
	}

	return head
}
