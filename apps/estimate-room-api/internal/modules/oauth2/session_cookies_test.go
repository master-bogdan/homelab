package oauth2

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionCookie_UsesTrustedForwardedProtoForSecureFlag(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/api/v1/auth/login", nil)
	req.Header.Set("X-Forwarded-Proto", "https")

	if !SessionCookie("session-123", req, true).Secure {
		t.Fatalf("expected session cookie to be Secure when trusted proxy headers indicate https")
	}
	if SessionCookie("session-123", req, false).Secure {
		t.Fatalf("expected session cookie to ignore forwarded proto when proxy headers are not trusted")
	}
}

func TestExpiredCookies_UseTrustedForwardedProtoForSecureFlag(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/api/v1/auth/logout", nil)
	req.Header.Set("X-Forwarded-Proto", "https")

	for _, cookie := range []*http.Cookie{
		ExpiredSessionCookie(req, true),
		ExpiredAccessTokenCookie(req, true),
		ExpiredRefreshTokenCookie(req, true),
	} {
		if !cookie.Secure {
			t.Fatalf("expected %s cookie to be Secure", cookie.Name)
		}
		if cookie.MaxAge != -1 {
			t.Fatalf("expected %s cookie to be expired", cookie.Name)
		}
	}
}
