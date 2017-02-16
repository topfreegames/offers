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

var _ = Describe("Game Handler", func() {
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	Describe("POST /game", func() {
		It("should return status code of 200", func() {
			gameReader := JSONFor(JSON{
				"Name":     uuid.NewV4().String(),
				"BundleID": "com.topfreegames.example",
			})
			request, _ := http.NewRequest("POST", "/game/upsert", gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(200))
		})

		It("should return status code of 400 if missing parameter", func() {
			gameReader := JSONFor(JSON{
				"Name": uuid.NewV4().String(),
			})
			request, _ := http.NewRequest("POST", "/game/upsert", gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(400))
		})

		It("should return status code of 400 if invalid name", func() {
			reallyBigName := "1234567890"
			for i := 0; i < 5; i++ {
				reallyBigName += reallyBigName
			}

			gameReader := JSONFor(JSON{
				"Name":     reallyBigName,
				"BundleID": "com.topfreegames.example",
			})
			request, _ := http.NewRequest("POST", "/game/upsert", gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(400))
			Expect(recorder.Body.String()).To(ContainSubstring("Payload is invalid: Name:"))
		})

		It("should return status code of 400 if invalid bundle id", func() {
			reallyBigName := "1234567890"
			for i := 0; i < 5; i++ {
				reallyBigName += reallyBigName
			}

			gameReader := JSONFor(JSON{
				"Name":     uuid.NewV4().String(),
				"BundleID": reallyBigName,
			})
			request, _ := http.NewRequest("POST", "/game/upsert", gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(400))
			Expect(recorder.Body.String()).To(ContainSubstring("Payload is invalid: BundleID:"))
		})

	})
})
