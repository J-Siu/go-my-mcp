package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPresentedKey(t *testing.T) {
	cases := []struct {
		name    string
		headers map[string]string
		want    string
		wantOk  bool
	}{
		{"no headers", nil, "", false},
		{"bearer", map[string]string{"Authorization": "Bearer secret"}, "secret", true},
		{"bearer lower case scheme", map[string]string{"Authorization": "bearer secret"}, "secret", true},
		{"x-api-key", map[string]string{APIKeyHeader: "secret"}, "secret", true},
		{"authorization wins", map[string]string{"Authorization": "Bearer a", APIKeyHeader: "b"}, "a", true},
		{"wrong scheme", map[string]string{"Authorization": "Basic secret"}, "", false},
		{"no scheme", map[string]string{"Authorization": "secret"}, "", false},
		{"empty bearer", map[string]string{"Authorization": "Bearer "}, "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/mcp", nil)
			for k, v := range c.headers {
				r.Header.Set(k, v)
			}
			got, ok := presentedKey(r)
			if got != c.want || ok != c.wantOk {
				t.Fatalf("presentedKey() = %q, %v; want %q, %v", got, ok, c.want, c.wantOk)
			}
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	cases := []struct {
		name   string
		auth   *TypeAuth
		header map[string]string
		want   int
	}{
		{"disabled without a key", nil, nil, http.StatusNoContent},
		{"disabled when key empty", &TypeAuth{Key: ""}, nil, http.StatusNoContent},
		{"missing key", &TypeAuth{Key: "secret"}, nil, http.StatusUnauthorized},
		{"wrong key", &TypeAuth{Key: "secret"}, map[string]string{"Authorization": "Bearer nope"}, http.StatusUnauthorized},
		{"right bearer", &TypeAuth{Key: "secret"}, map[string]string{"Authorization": "Bearer secret"}, http.StatusNoContent},
		{"right x-api-key", &TypeAuth{Key: "secret"}, map[string]string{APIKeyHeader: "secret"}, http.StatusNoContent},
		{"prefix of key", &TypeAuth{Key: "secret"}, map[string]string{"Authorization": "Bearer secre"}, http.StatusUnauthorized},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/mcp", nil)
			for k, v := range c.header {
				r.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			c.auth.Middleware(okHandler).ServeHTTP(w, r)
			if w.Code != c.want {
				t.Fatalf("status = %d; want %d", w.Code, c.want)
			}
			if c.want == http.StatusUnauthorized && w.Header().Get("WWW-Authenticate") == "" {
				t.Fatal("missing WWW-Authenticate on 401")
			}
		})
	}
}
