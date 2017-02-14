// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package testing

import (
	"github.com/go-testfixtures/testfixtures"
	runner "github.com/mgutz/dat/sqlx-runner"
	"github.com/topfreegames/offers/models"
)

var (
	fixtures *testfixtures.Context
)

//GetTestDB returns a connection to the test database
func GetTestDB() (*runner.DB, error) {
	return models.GetDB(
		"localhost", "offers_test", 8585, "disable",
		"offers_test", "",
		10, 10,
	)
}

//LoadFixtures into the DB
func LoadFixtures(db *runner.DB) error {
	var err error

	if fixtures == nil {
		// creating the context that hold the fixtures
		// see about all compatible databases in this page below
		fixtures, err = testfixtures.NewFolder(db.DB.DB, &testfixtures.PostgreSQL{}, "../fixtures")
		if err != nil {
			return err
		}
	}

	if err := fixtures.Load(); err != nil {
		return err
	}

	return nil
}
