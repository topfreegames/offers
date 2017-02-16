// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

//MixedMetricsReporter calls other metrics reporters
type MixedMetricsReporter struct {
	MetricsReporters []MetricsReporter
}

//NewMixedMetricsReporter ctor
func NewMixedMetricsReporter() *MixedMetricsReporter {
	return &MixedMetricsReporter{
		MetricsReporters: []MetricsReporter{},
	}
}

//WithSegment that calls all the other metrics reporters
func (m *MixedMetricsReporter) WithSegment(name string, f func() error) error {
	if m == nil {
		return f()
	}

	ff := f
	for _, mr := range m.MetricsReporters {
		ff = func() error {
			return mr.WithSegment(name, ff)
		}
	}

	return ff()
}

//AddReporter to metrics reporter
func (m *MixedMetricsReporter) AddReporter(mr MetricsReporter) {
	m.MetricsReporters = append(m.MetricsReporters, mr)
}
