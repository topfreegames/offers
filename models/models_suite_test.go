// offers api
// https://github.com/topfreegames/offers api
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	runner "github.com/mgutz/dat/sqlx-runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	oTesting "github.com/topfreegames/offers/testing"
)

var db runner.Connection

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Offers API - Models Suite")
}

var _ = BeforeSuite(func() {
	var err error
	db, err = oTesting.GetTestDB()
	Expect(err).NotTo(HaveOccurred())

	err = oTesting.LoadFixtures(db.(*runner.DB))
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := db.(*runner.DB).DB.Close()
	Expect(err).NotTo(HaveOccurred())
})
