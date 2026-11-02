/*
Copyright © 2026 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package server

import (
	"crypto/subtle"
	"log/slog"
	"net/http"
	"strings"
)

// APIKeyHeader - alternative to the standard Authorization header
const APIKeyHeader = "X-API-Key"

// TypeAuth - API key check applied to every request before the MCP handler
type TypeAuth struct {
	Key string
}

// Middleware - reject requests that do not carry the configured API key.
//
// A nil TypeAuth, or one without a key, does no checking, so access control
// stays off unless a key is configured.
func (t *TypeAuth) Middleware(next http.Handler) http.Handler {
	if t == nil || t.Key == "" {
		return next
	}
	expected := []byte(t.Key)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		presented, ok := presentedKey(r)
		// constant time compare, so the key cannot be recovered by timing
		if !ok || subtle.ConstantTimeCompare([]byte(presented), expected) != 1 {
			slog.Warn("TypeAuth.Middleware: unauthorized",
				"remote", r.RemoteAddr, "path", r.URL.Path)
			w.Header().Set("WWW-Authenticate", `Bearer realm="go-my-mcp"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Extract the key from either accepted header. `Authorization: Bearer` wins
// when both are present.
func presentedKey(r *http.Request) (string, bool) {
	if auth := r.Header.Get("Authorization"); auth != "" {
		scheme, token, found := strings.Cut(auth, " ")
		if !found || !strings.EqualFold(scheme, "Bearer") || token == "" {
			// present but malformed; do not fall through to the other header,
			// so a broken Authorization is reported rather than silently ignored
			return "", false
		}
		return strings.TrimSpace(token), true
	}
	if key := r.Header.Get(APIKeyHeader); key != "" {
		return key, true
	}
	return "", false
}
