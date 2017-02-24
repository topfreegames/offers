// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
				Expect(recorder.Body.String()).To(Equal(`{"healthy": true}`))
			})

			It("returns the version as a header", func() {
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Header().Get("X-Offers-Version")).To(Equal("0.1.0"))
			})

			It("returns status code of 500 if database is unavailable", func() {
				oldDB := app.DB
				db, err := GetTestDB()
				Expect(err).NotTo(HaveOccurred())
				app.DB = db
				app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
				app.Router.ServeHTTP(recorder, request)

				Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
				var obj map[string]interface{}
				err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
				Expect(err).NotTo(HaveOccurred())
				Expect(obj["code"]).To(Equal("OFF-000"))
				Expect(obj["error"]).To(Equal("DatabaseError"))
				Expect(obj["description"]).To(Equal("sql: database is closed"))
				app.DB = oldDB // avoid errors in after each
			})
		})
	})
})
