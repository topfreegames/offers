// offers api
// https://github.com/topfree/ames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/topfreegames/offers/testing"
)

var _ = Describe("Offer Template Handler", func() {
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	Describe("PUT /offer-template", func() {
		It("should return status code 200 for valid parameters", func() {
			offerTemplateReader := JSONFor(JSON{
				"ID":        "56fc0477-39f1-485c-898e-4909e9155eb1",
				"Name":      "New Awesome Game",
				"ProductID": "com.tfg.example",
				"GameID":    "nonexisting-game-id",
				"Contents":  "{\"gems\": 5, \"gold\": 100}",
				"Period":    "{\"type\": \"once\"}",
				"Frequency": "{\"every\": 24, \"unit\": \"hour\"}",
				"Trigger":   "{\"from\": 1487280506875, \"to\": 1487366964730}",
			})

			request, _ := http.NewRequest("PUT", "/offer-templates", offerTemplateReader)
			app.Router.ServeHTTP(recorder, request)
			//Expect(recorder.Body).To(Equal("qwe"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})
	})
})
