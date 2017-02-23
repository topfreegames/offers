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

var _ = Describe("Trigger Time", func() {
	times := models.Times{
		From: 5,
		To:   10,
	}
	Describe("IsTriggered", func() {
		It("should return true if from <= now <= to", func() {
			now := time.Unix(7, 0)
			trigger := models.TimeTrigger{}
			isTriggered := trigger.IsTriggered(times, now)
			Expect(isTriggered).To(BeTrue())
		})

		It("should return false if now < from", func() {
			now := time.Unix(3, 0)
			trigger := models.TimeTrigger{}
			isTriggered := trigger.IsTriggered(times, now)
			Expect(isTriggered).To(BeFalse())
		})

		It("should return false if now > to", func() {
			now := time.Unix(13, 0)
			trigger := models.TimeTrigger{}
			isTriggered := trigger.IsTriggered(times, now)
			Expect(isTriggered).To(BeFalse())
		})
	})
})
