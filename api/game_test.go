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
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("Game Handler", func() {
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	Describe("POST /game", func() {
		It("should return status code of 200", func() {
			game := `{"Name": "` + uuid.NewV4().String() + `", "BundleID": "com.topfreegames.example"}`
			request, _ := http.NewRequest("POST", "/game/upsert", strings.NewReader(game))

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(200))
		})

		It("should return status code of 400 if missing parameter", func() {
			game := `{"Name": "` + uuid.NewV4().String() + `"}`
			request, _ := http.NewRequest("POST", "/game/upsert", strings.NewReader(game))

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(400))
		})
	})
})
