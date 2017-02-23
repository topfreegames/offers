// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topfreegames/offers/models"
)

var _ = Describe("Trigger Default", func() {
	Describe("IsTriggered", func() {
		It("should return true", func() {
			trigger := models.DefaultTrigger{}
			isTriggered := trigger.IsTriggered(nil, nil)
			Expect(isTriggered).To(BeTrue())
		})
	})
})
