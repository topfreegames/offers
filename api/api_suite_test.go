// offers api
// https://github.com/topfreegames/offers api
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	"io"
	"time"

	"github.com/Sirupsen/logrus"
	runner "github.com/mgutz/dat/sqlx-runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/topfreegames/offers/api"
	"github.com/topfreegames/offers/models"
	oTesting "github.com/topfreegames/offers/testing"
)

var app *api.App
var db runner.Connection
var closer io.Closer

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
	Expect(db).NotTo(BeNil())

	models.ShouldPing(db.(*runner.DB).DB.DB, 50*time.Millisecond)

	err = oTesting.LoadFixtures(db.(*runner.DB))
	Expect(err).NotTo(HaveOccurred())

	config, err := oTesting.GetDefaultConfig()
	Expect(err).NotTo(HaveOccurred())
	app, err = api.NewApp("0.0.0.0", 8889, config, false, l, nil)
	Expect(err).NotTo(HaveOccurred())
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
