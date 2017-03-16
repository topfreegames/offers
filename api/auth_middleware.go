// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"context"
	"net/http"
)

//AuthMiddleware automatically adds a user email to the context
type AuthMiddleware struct {
	App          *App
	Next         http.Handler
	useBasicAuth bool
}

const userEmailKey = contextKey("userEmail")

func newContextWithUserEmail(ctx context.Context, r *http.Request) context.Context {
	userEmail := r.Header.Get("x-forwarded-email")
	c := context.WithValue(ctx, userEmailKey, userEmail)
	return c
}

func userEmailFromContext(ctx context.Context) string {
	return ctx.Value(userEmailKey).(string)
}

// ServeHTTP method
func (m *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContextWithUserEmail(r.Context(), r)

	basicAuthUser := m.App.Config.GetString("basicauth.username")
	basicAuthPass := m.App.Config.GetString("basicauth.password")
	if m.useBasicAuth && basicAuthUser != "" && basicAuthPass != "" {
		user, pass, ok := r.BasicAuth()
		if !ok {
			Write(w, http.StatusUnauthorized, "Authentication failed.")
			return
		} else if user != basicAuthUser || pass != basicAuthPass {
			Write(w, http.StatusUnauthorized, "Authentication failed.")
			return
		}
	}
	// Call the next middleware/handler in chain
	m.Next.ServeHTTP(w, r.WithContext(ctx))
}

//SetNext middleware
func (m *AuthMiddleware) SetNext(next http.Handler) {
	m.Next = next
}
