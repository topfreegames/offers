// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import "time"

//MetricsReporter is a contract for reporters of metrics
type MetricsReporter interface {
	StartSegment(string) map[string]interface{}
	EndSegment(map[string]interface{}, string)

	StartDatastoreSegment(datastore, collection, operation string) map[string]interface{}
	EndDatastoreSegment(map[string]interface{})

	StartExternalSegment(string) map[string]interface{}
	EndExternalSegment(map[string]interface{})
}

//Clock returns the time
type Clock interface {
	GetTime() time.Time
}

//Trigger return true if offer is triggered
type Trigger interface {
  IsTriggered(interface{}, interface{}) bool
}
