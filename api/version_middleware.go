// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import "net/http"

// VersionMiddleware adds the version to the request
type VersionMiddleware struct {
	next http.Handler
}

//ServeHTTP method
func (m *VersionMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Offers-Version", "0.1.0")
	m.next.ServeHTTP(w, r)
}

//SetNext handler
func (m *VersionMiddleware) SetNext(next http.Handler) {
	m.next = next
}
