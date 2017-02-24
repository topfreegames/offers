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

	newrelic "github.com/newrelic/go-agent"
)

//NewRelicMiddleware handles logging
type NewRelicMiddleware struct {
	App  *App
	Next http.Handler
}

const newRelicTransactionKey string = "newRelicTransaction"

func newContextWithNewRelicTransaction(ctx context.Context, txn newrelic.Transaction, r *http.Request) context.Context {
	c := context.WithValue(ctx, newRelicTransactionKey, txn)
	return c
}

func (m *NewRelicMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if m.App.NewRelic != nil {
		txn := m.App.NewRelic.StartTransaction(fmt.Sprintf("%s %s", r.Method, r.URL.Path), w, r)
		defer txn.End()
		ctx = newContextWithNewRelicTransaction(r.Context(), txn, r)

		mr := metricsReporterFromCtx(ctx)
		if mr != nil {
			mr.AddReporter(&NewRelicMetricsReporter{
				App:         m.App,
				Transaction: txn,
			})
		}
	}

	// Call the next middleware/handler in chain
	m.Next.ServeHTTP(w, r.WithContext(ctx))
}

//SetNext middleware
func (m *NewRelicMiddleware) SetNext(next http.Handler) {
	m.Next = next
}
