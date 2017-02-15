// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

//MetricsReporter is a contract for reporters of metrics
type MetricsReporter interface {
	WithSegment(string, func() error) error
}
