// offers api
// https://github.com/topfreegames/offers api
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"

	"testing"

	"github.com/topfreegames/offers/api"
	oTesting "github.com/topfreegames/offers/testing"
)

var app *api.App
var db runner.Connection
var closer io.Closer
var config *viper.Viper

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Offers API - API Suite")
}

var _ = BeforeSuite(func() {
	l := logrus.New()
	l.Level = logrus.FatalLevel

	var err error
	db, err = oTesting.GetTestDB()
	Expect(err).NotTo(HaveOccurred())

	err = oTesting.ClearOfferPlayers(db)
	Expect(err).NotTo(HaveOccurred())

	err = oTesting.LoadFixtures(db)
	Expect(err).NotTo(HaveOccurred())

	config, err = oTesting.GetDefaultConfig()
	Expect(err).NotTo(HaveOccurred())
	clock := oTesting.MockClock{
		CurrentTime: 1486678000,
	}
	app, err = api.NewApp("0.0.0.0", 8889, config, false, l, clock)
	Expect(err).NotTo(HaveOccurred())
})

var _ = BeforeEach(func() {
	tx, err := db.Begin()
	Expect(err).NotTo(HaveOccurred())
	app.DB = tx
})

var _ = AfterEach(func() {
	if !app.DB.(*runner.Tx).IsRollbacked {
		err := app.DB.(*runner.Tx).Rollback()
		Expect(err).NotTo(HaveOccurred())
	}
	app.DB = db
	app.Clock = oTesting.MockClock{
		CurrentTime: 1486678000,
	}
	app.Cache.Flush()
})

var _ = AfterSuite(func() {
	if db != nil {
		err := db.(*runner.DB).DB.Close()
		Expect(err).NotTo(HaveOccurred())
		db = nil
	}

	if closer != nil {
		closer.Close()
		closer = nil
	}
})
