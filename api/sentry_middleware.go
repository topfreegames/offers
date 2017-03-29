// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"net/http"

	raven "github.com/getsentry/raven-go"
)

// SentryMiddleware adds the version to the request
type SentryMiddleware struct {
	next http.Handler
}

//ServeHTTP method
func (m *SentryMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	recoveryHandler := raven.RecoveryHandler(m.next.ServeHTTP)
	recoveryHandler(w, r)
}

//SetNext handler
func (m *SentryMiddleware) SetNext(next http.Handler) {
	m.next = next
}
