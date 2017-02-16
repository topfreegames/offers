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
	uuid "github.com/satori/go.uuid"
)

//LoggingMiddleware handles logging
type LoggingMiddleware struct {
	App  *App
	Next http.Handler
}

const requestIDKey string = "requestID"
const loggerKey string = "logger"

func newContextWithRequestIDAndLogger(logger logrus.FieldLogger, ctx context.Context, r *http.Request) context.Context {
	reqID := uuid.NewV4().String()
	l := logger.WithField("requestID", reqID)

	c := context.WithValue(ctx, requestIDKey, reqID)
	c = context.WithValue(c, loggerKey, l)
	return c
}

func requestIDFromContext(ctx context.Context) string {
	return ctx.Value(requestIDKey).(string)
}

func loggerFromContext(ctx context.Context) logrus.FieldLogger {
	return ctx.Value(loggerKey).(logrus.FieldLogger)
}

func (m *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContextWithRequestIDAndLogger(m.App.Logger, r.Context(), r)

	start := time.Now()
	defer func() {
		l := loggerFromContext(ctx)
		l.WithFields(logrus.Fields{
			"path":            r.URL.Path,
			"requestDuration": time.Since(start).Nanoseconds(),
		}).Info("Request completed.")
	}()

	// Call the next middleware/handler in chain
	m.Next.ServeHTTP(w, r.WithContext(ctx))
}

//SetNext middleware
func (m *LoggingMiddleware) SetNext(next http.Handler) {
	m.Next = next
}