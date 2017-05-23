// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"net/http"

	"github.com/DataDog/dd-trace-go/tracer"
	"github.com/gorilla/mux"
)

//DataDogMiddleware handles logging
type DataDogMiddleware struct {
	App  *App
	Next http.Handler
}

func (m *DataDogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, _ := mux.CurrentRoute(r).GetPathTemplate()
	span := tracer.NewRootSpan("http.request", "offers-api", route)
	defer span.Finish()
	ctx := span.Context(r.Context())
	// Call the next middleware/handler in chain
	m.Next.ServeHTTP(w, r.WithContext(ctx))
}

//SetNext middleware
func (m *DataDogMiddleware) SetNext(next http.Handler) {
	m.Next = next
}
