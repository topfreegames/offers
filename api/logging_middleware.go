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
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

//LoggingMiddleware handles logging
type LoggingMiddleware struct {
	App  *App
	Next http.Handler
}

const requestIDKey = contextKey("requestID")
const loggerKey = contextKey("logger")

func newContextWithRequestIDAndLogger(ctx context.Context, logger logrus.FieldLogger, r *http.Request) context.Context {
	reqID := uuid.NewV4().String()
	l := logger.WithField("requestID", reqID)

	c := context.WithValue(ctx, requestIDKey, reqID)
	c = context.WithValue(c, loggerKey, l)
	return c
}

func loggerFromContext(ctx context.Context) logrus.FieldLogger {
	return ctx.Value(loggerKey).(logrus.FieldLogger)
}

// ServeHTTP method
func (m *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContextWithRequestIDAndLogger(r.Context(), m.App.Logger, r)

	start := time.Now()
	defer func() {
		l := loggerFromContext(ctx)
		status := getStatusFromResponseWriter(w)
		route, _ := mux.CurrentRoute(r).GetPathTemplate()
		// request failed
		if status > 399 && status < 500 {
			l.WithFields(logrus.Fields{
				"path":            r.URL.Path,
				"route":           route,
				"requestDuration": time.Since(start).Nanoseconds(),
				"status":          status,
			}).Warn("Request failed.")
		} else if status > 499 { // request is ok, but server failed
			l.WithFields(logrus.Fields{
				"path":            r.URL.Path,
				"route":           route,
				"requestDuration": time.Since(start).Nanoseconds(),
				"status":          status,
			}).Error("Response failed.")
		} else { // Everything went ok
			l.WithFields(logrus.Fields{
				"path":            r.URL.Path,
				"route":           route,
				"requestDuration": time.Since(start).Nanoseconds(),
				"status":          status,
			}).Info("Request successful.")
		}
	}()

	// Call the next middleware/handler in chain
	m.Next.ServeHTTP(w, r.WithContext(ctx))
}

//SetNext middleware
func (m *LoggingMiddleware) SetNext(next http.Handler) {
	m.Next = next
}
