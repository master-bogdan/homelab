package oauth2

import "net/http"

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

func SessionCookie(sessionID string, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		HttpOnly: true,
		Secure:   secure,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
}

func ExpiredSessionCookie(secure bool) *http.Cookie {
	return expiredCookie(SessionCookieName, secure)
}

func ExpiredAccessTokenCookie(secure bool) *http.Cookie {
	return expiredCookie(AccessTokenCookieName, secure)
}

func ExpiredRefreshTokenCookie(secure bool) *http.Cookie {
	return expiredCookie(RefreshTokenCookieName, secure)
}

func expiredCookie(name string, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		HttpOnly: true,
		Secure:   secure,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
}
