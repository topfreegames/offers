// offers api
// https://github.com/topfreegames/offers api
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"

	"testing"

	oTesting "github.com/topfreegames/offers/testing"
	"github.com/topfreegames/offers/util"
)

var conn runner.Connection
var db *runner.Tx
var redisClient *util.RedisClient

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

	redisClient, err = oTesting.GetTestRedis()
	Expect(err).NotTo(HaveOccurred())
})

var _ = BeforeEach(func() {
	var err error
	db, err = conn.Begin()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterEach(func() {
	err := db.Rollback()
	Expect(err).NotTo(HaveOccurred())
	db = nil
	status := redisClient.Client.FlushAll()
	Expect(status.Err()).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := conn.(*runner.DB).DB.Close()
	Expect(err).NotTo(HaveOccurred())
})
