package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/master-bogdan/estimate-room-api/internal/modules/oauth2"
)

func TestOauth2SessionCookie_UsesTrustedForwardedProtoForSecureFlag(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/api/v1/auth/login", nil)
	req.Header.Set("X-Forwarded-Proto", "https")

	if !oauth2.Oauth2SessionCookie("session-123", req, true).Secure {
		t.Fatalf("expected session cookie to be Secure when trusted proxy headers indicate https")
	}
	if oauth2.Oauth2SessionCookie("session-123", req, false).Secure {
		t.Fatalf("expected session cookie to ignore forwarded proto when proxy headers are not trusted")
	}
}

func TestExpiredOauth2Cookies_UseTrustedForwardedProtoForSecureFlag(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/api/v1/auth/logout", nil)
	req.Header.Set("X-Forwarded-Proto", "https")

	for _, cookie := range []*http.Cookie{
		oauth2.ExpiredOauth2SessionCookie(req, true),
		oauth2.ExpiredOauth2AccessTokenCookie(req, true),
		oauth2.ExpiredOauth2RefreshTokenCookie(req, true),
	} {
		if !cookie.Secure {
			t.Fatalf("expected %s cookie to be Secure", cookie.Name)
		}
		if cookie.MaxAge != -1 {
			t.Fatalf("expected %s cookie to be expired", cookie.Name)
		}
	}
}
