// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import newrelic "github.com/newrelic/go-agent"

//NewRelicMetricsReporter reports metrics to new relic
type NewRelicMetricsReporter struct {
	App         *App
	Transaction newrelic.Transaction
}

//WithSegment starts a segment in the transaction
func (r *NewRelicMetricsReporter) WithSegment(name string, f func() error) error {
	if r.Transaction == nil {
		return f()
	}

	defer newrelic.StartSegment(r.Transaction, name).End()
	return f()
}
