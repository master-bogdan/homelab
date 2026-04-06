package oauth2

import (
	"net/http"
	"strings"
)

const (
	Oauth2SessionCookieName      = "session_id"
	Oauth2AccessTokenCookieName  = "access_token"
	Oauth2RefreshTokenCookieName = "refresh_token"
)

func ReadOauth2SessionID(r *http.Request) string {
	if r == nil {
		return ""
	}

	if cookie, err := r.Cookie(Oauth2SessionCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	if header := r.Header.Get("X-Session-Id"); header != "" {
		return header
	}

	return ""
}

func Oauth2SessionCookie(sessionID string, r *http.Request, trustProxyHeaders bool) *http.Cookie {
	return cookieWithSettings(Oauth2SessionCookieName, sessionID, 0, isSecureRequest(r, trustProxyHeaders))
}

func Oauth2AccessTokenCookie(token string, r *http.Request, trustProxyHeaders bool) *http.Cookie {
	return cookieWithSettings(
		Oauth2AccessTokenCookieName,
		token,
		int(accessTokenTTL.Seconds()),
		isSecureRequest(r, trustProxyHeaders),
	)
}

func Oauth2RefreshTokenCookie(token string, r *http.Request, trustProxyHeaders bool) *http.Cookie {
	return cookieWithSettings(
		Oauth2RefreshTokenCookieName,
		token,
		int(refreshTokenTTL.Seconds()),
		isSecureRequest(r, trustProxyHeaders),
	)
}

func ExpiredOauth2SessionCookie(r *http.Request, trustProxyHeaders bool) *http.Cookie {
	return expiredCookie(Oauth2SessionCookieName, isSecureRequest(r, trustProxyHeaders))
}

func ExpiredOauth2AccessTokenCookie(r *http.Request, trustProxyHeaders bool) *http.Cookie {
	return expiredCookie(Oauth2AccessTokenCookieName, isSecureRequest(r, trustProxyHeaders))
}

func ExpiredOauth2RefreshTokenCookie(r *http.Request, trustProxyHeaders bool) *http.Cookie {
	return expiredCookie(Oauth2RefreshTokenCookieName, isSecureRequest(r, trustProxyHeaders))
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
