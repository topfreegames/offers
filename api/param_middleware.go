// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/topfreegames/extensions/middleware"
	"github.com/topfreegames/offers/errors"
)

//ParamMiddleware add into the model the parameters that came in the URI
type ParamMiddleware struct {
	App       *App
	Validator func(string) bool
	Next      http.Handler
}

const paramKey = contextKey("param")

//NewParamKeyMiddleware constructs a new param key middleware
func NewParamKeyMiddleware(app *App, f func(string) bool) *ParamMiddleware {
	m := &ParamMiddleware{App: app, Validator: f}
	return m
}

func newContextWithParamKey(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, paramKey, id)
}

func paramKeyFromContext(ctx context.Context) string {
	param := ctx.Value(paramKey)
	if param == nil {
		return ""
	}
	return param.(string)
}

func (m *ParamMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	valid := m.Validator(id)

	if !valid {
		err := fmt.Errorf("ID: " + id + " does not validate;")
		l := middleware.GetLogger(r.Context())
		l.WithError(err).Error("Payload could not be decoded.")
		vErr := errors.NewValidationFailedError(err)
		WriteBytes(w, http.StatusUnprocessableEntity, vErr.Serialize())
		return
	}

	ctx := newContextWithParamKey(r.Context(), id)
	m.Next.ServeHTTP(w, r.WithContext(ctx))
}

//SetNext handler
func (m *ParamMiddleware) SetNext(next http.Handler) {
	m.Next = next
}
