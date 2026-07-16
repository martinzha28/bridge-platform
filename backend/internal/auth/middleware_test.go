package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestOptionalMiddlewarePassesGuestThrough(t *testing.T) {
	var sawUser bool
	h := OptionalMiddleware("secret")(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		_, sawUser = UserIDFromContext(r.Context())
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ws", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if sawUser {
		t.Fatal("guest request should carry no user ID")
	}
}

func TestOptionalMiddlewareAttachesValidUser(t *testing.T) {
	const secret = "secret"
	id := uuid.New()
	token, err := IssueToken(id, secret)
	if err != nil {
		t.Fatalf("IssueToken: %v", err)
	}

	var got uuid.UUID
	var ok bool
	h := OptionalMiddleware(secret)(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		got, ok = UserIDFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.AddCookie(&http.Cookie{Name: cookieName, Value: token})
	h.ServeHTTP(httptest.NewRecorder(), req)

	if !ok || got != id {
		t.Fatalf("user = %v (ok=%v), want %v", got, ok, id)
	}
}

func TestOptionalMiddlewareIgnoresBadToken(t *testing.T) {
	var sawUser bool
	h := OptionalMiddleware("secret")(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		_, sawUser = UserIDFromContext(r.Context())
	}))

	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.AddCookie(&http.Cookie{Name: cookieName, Value: "not-a-jwt"})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK || sawUser {
		t.Fatalf("bad token should pass through unauthenticated (code=%d, sawUser=%v)", rec.Code, sawUser)
	}
}
