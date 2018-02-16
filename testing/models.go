// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package testing

import (
	"github.com/go-testfixtures/testfixtures"
	"github.com/topfreegames/offers/models"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"
)

var (
	fixtures *testfixtures.Context
)

//GetTestDB returns a connection to the test database
func GetTestDB() (runner.Connection, error) {
	return models.GetDB(
		"localhost", "offers_test", 8585, "disable",
		"offers_test", "",
		10, 10, 100,
	)
}

//GetPerfDB returns a connection to the perf database
func GetPerfDB() (runner.Connection, error) {
	return models.GetDB(
		"localhost", "offers_perf", 8585, "disable",
		"offers_perf", "",
		10, 10, 100,
	)
}

//ClearOfferPlayers
func ClearOfferPlayers(db runner.Connection) error {
	_, err := db.DeleteFrom("offer_players").Exec()
	return err
}

//LoadFixtures into the DB
func LoadFixtures(db runner.Connection) error {
	var err error

	conn := db.(*runner.DB).DB.DB

	if fixtures == nil {
		// creating the context that hold the fixtures
		// see about all compatible databases in this page below
		fixtures, err = testfixtures.NewFolder(conn, &testfixtures.PostgreSQL{}, "../fixtures")
		if err != nil {
			return err
		}
	}

	if err := fixtures.Load(); err != nil {
		return err
	}

	return nil
}

//FakeMetricsReporter is a fake metric reporter for testing
type FakeMetricsReporter struct{}

//StartSegment mocks a metric reporter
func (mr FakeMetricsReporter) StartSegment(key string) map[string]interface{} {
	return map[string]interface{}{}
}

//EndSegment mocks a metric reporter
func (mr FakeMetricsReporter) EndSegment(m map[string]interface{}, key string) {}

//StartDatastoreSegment mocks a metric reporter
func (mr FakeMetricsReporter) StartDatastoreSegment(datastore, collection, operation string) map[string]interface{} {
	return map[string]interface{}{}
}

//EndDatastoreSegment mocks a metric reporter
func (mr FakeMetricsReporter) EndDatastoreSegment(m map[string]interface{}) {}

//StartExternalSegment mocks a metric reporter
func (mr FakeMetricsReporter) StartExternalSegment(key string) map[string]interface{} {
	return map[string]interface{}{}
}

//EndExternalSegment mocks a metric reporter
func (mr FakeMetricsReporter) EndExternalSegment(m map[string]interface{}) {}
