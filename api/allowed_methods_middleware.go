// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"net/http"
	"strings"
)

// AllowedMethodsMiddleware ensures that only a url can me requested with
// a specific method, else returns a 400 Bad Request
type AllowedMethodsMiddleware struct {
	AllowedMethods []string
	next           http.Handler
}

//AllowedMethods constructs a middleware
func AllowedMethods(methods ...string) Middleware {
	return &AllowedMethodsMiddleware{
		AllowedMethods: methods,
	}
}

//ServeHTTP method
func (m *AllowedMethodsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, method := range m.AllowedMethods {
		if strings.ToLower(r.Method) == strings.ToLower(method) {
			m.next.ServeHTTP(w, r)
			return
		}
	}

	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

//SetNext handler
func (m *AllowedMethodsMiddleware) SetNext(next http.Handler) {
	m.next = next
}
