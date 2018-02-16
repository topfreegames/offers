// offers api
// https://github.com/topfreegames/offers api
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"

	"testing"

	"github.com/pmylund/go-cache"
	oTesting "github.com/topfreegames/offers/testing"
)

var conn runner.Connection
var db *runner.Tx
var offersCache *cache.Cache

func TestApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Offers API - Models Suite")
}

var _ = BeforeSuite(func() {
	var err error
	conn, err = oTesting.GetTestDB()
	Expect(err).NotTo(HaveOccurred())

	err = oTesting.LoadFixtures(conn)
	Expect(err).NotTo(HaveOccurred())

	offersCache = cache.New(300*time.Second, 30*time.Second)
})

var _ = BeforeEach(func() {
	var err error
	db, err = conn.Begin()
	Expect(err).NotTo(HaveOccurred())
	offersCache.Flush()
})

var _ = AfterEach(func() {
	err := db.Rollback()
	Expect(err).NotTo(HaveOccurred())
	db = nil
})

var _ = AfterSuite(func() {
	err := conn.(*runner.DB).DB.Close()
	Expect(err).NotTo(HaveOccurred())
})
