// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topfreegames/offers/models"
)

var _ = Describe("Clock", func() {
	Describe("GetTime", func() {
		It("should return current time", func() {
			clock := models.RealClock{}
			t := clock.GetTime()
			Expect(time.Now().Unix() - t.Unix()).To(BeNumerically("<", 5))
		})
	})
})
