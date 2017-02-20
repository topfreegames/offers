// offers api
// https://github.com/topfreegames/offers
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
	uuid "github.com/satori/go.uuid"
	. "github.com/topfreegames/offers/testing"
)

var _ = Describe("Healthcheck Handler", func() {
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		// Record HTTP responses.
		recorder = httptest.NewRecorder()
	})

	Describe("GET /healthcheck", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/healthcheck", nil)
		})

		Context("when all services healthy", func() {
			It("returns a status code of 200", func() {
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns working string", func() {
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Body.String()).To(Equal("WORKING"))
			})

			It("returns the version as a header", func() {
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Header().Get("X-Offers-Version")).To(Equal("0.1.0"))
			})
		})
	})

	Describe("PUT /offer", func() {
		It("should insert valid offer", func() {
			//Given
			offerReader := JSONFor(JSON{
				"ID":              uuid.NewV4().String(),
				"GameID":          "offers-game",
				"OfferTemplateID": "dd21ec96-2890-4ba0-b8e2-40ea67196990",
				"PlayerID":        "player-1",
			})
			request, _ := http.NewRequest("PUT", "/offer", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

		It("should receive status code 400 if OfferTemplateID not exists", func() {
			//Given
			offerReader := JSONFor(JSON{
				"ID":              uuid.NewV4().String(),
				"GameID":          "offers-game",
				"OfferTemplateID": uuid.NewV4().String(),
				"PlayerID":        "player-1",
			})
			request, _ := http.NewRequest("PUT", "/offer", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		})

		It("should receive status code 400 if game-id doesn't exist", func() {
			//Given
			offerReader := JSONFor(JSON{
				"ID":              uuid.NewV4().String(),
				"GameID":          "non-existing-offers-game",
				"OfferTemplateID": "dd21ec96-2890-4ba0-b8e2-40ea67196990",
				"PlayerID":        "player-1",
			})
			request, _ := http.NewRequest("PUT", "/offer", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("GET /get-offers", func() {

	})

	Describe("PUT /claim-offer", func() {
		It("should claim valid offer", func() {
			//Given
			offerReader := JSONFor(JSON{
				"ID":     "56fc0477-39f1-485c-898e-4909e9155eb1",
				"GameID": "offers-game",
			})
			request, _ := http.NewRequest("PUT", "/claim-offer", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			//Expect(recorder.Code).To(Equal(http.StatusOK))
		})
	})
})
