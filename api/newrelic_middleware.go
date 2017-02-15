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

	newrelic "github.com/newrelic/go-agent"
)

//NewRelicMiddleware handles logging
type NewRelicMiddleware struct {
	App  *App
	Next http.Handler
}

const newRelicTransactionKey string = "newRelicTransaction"

func newContextWithNewRelicTransaction(txn newrelic.Transaction, ctx context.Context, r *http.Request) context.Context {
	c := context.WithValue(ctx, newRelicTransactionKey, txn)
	return c
}

func newrelicTransactionFromContext(ctx context.Context) newrelic.Transaction {
	return ctx.Value(newRelicTransactionKey).(newrelic.Transaction)
}

func (m *NewRelicMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	txn := m.App.NewRelic.StartTransaction(r.URL.Path, w, r)
	defer txn.End()
	ctx := newContextWithNewRelicTransaction(txn, r.Context(), r)

	// Call the next middleware/handler in chain
	m.Next.ServeHTTP(w, r.WithContext(ctx))
}

//SetNext middleware
func (m *NewRelicMiddleware) SetNext(next http.Handler) {
	m.Next = next
}
