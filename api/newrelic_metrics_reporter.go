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

//StartSegment starts segment
func (r *NewRelicMetricsReporter) StartSegment(name string) map[string]interface{} {
	if r.Transaction == nil {
		return nil
	}

	return map[string]interface{}{
		"segment": newrelic.StartSegment(r.Transaction, name),
	}
}

//EndSegment stops segment
func (r *NewRelicMetricsReporter) EndSegment(data map[string]interface{}, name string) {
	if r.Transaction == nil {
		return
	}

	data["segment"].(newrelic.Segment).End()
}
