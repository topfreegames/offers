// offers api
// https://github.com/topfree/ames/offers
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
	uuid "github.com/satori/go.uuid"
	. "github.com/topfreegames/offers/testing"
)

var _ = Describe("Game Handler", func() {
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	Describe("PUT /games", func() {
		It("should return status code of 200", func() {
			id := uuid.NewV4().String()
			gameReader := JSONFor(JSON{
				"ID":       id,
				"Name":     "Game Awesome Name",
				"BundleID": "com.topfreegames.example",
			})
			request, _ := http.NewRequest("PUT", "/games", gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["gameId"]).To(Equal(id))
		})

		It("should return status code of 422 if missing parameter", func() {
			gameReader := JSONFor(JSON{
				"Name": "Game Awesome Name",
			})
			request, _ := http.NewRequest("PUT", "/games", gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: non zero value required;BundleID: non zero value required;"))
		})

		It("should return status code of 422 if invalid name", func() {
			reallyBigName := "1234567890"
			for i := 0; i < 5; i++ {
				reallyBigName += reallyBigName
			}

			gameReader := JSONFor(JSON{
				"ID":       "game-id",
				"Name":     reallyBigName,
				"BundleID": "com.topfreegames.example",
			})
			request, _ := http.NewRequest("PUT", "/games", gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(BeEquivalentTo("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(ContainSubstring("does not validate as stringlength(1|255);"))
		})

		It("should return status code of 422 if invalid bundle id", func() {
			reallyBigName := "1234567890"
			for i := 0; i < 5; i++ {
				reallyBigName += reallyBigName
			}

			gameReader := JSONFor(JSON{
				"ID":       "game-id",
				"Name":     uuid.NewV4().String(),
				"BundleID": reallyBigName,
			})
			request, _ := http.NewRequest("PUT", "/games", gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))

			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(BeEquivalentTo("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(ContainSubstring("does not validate as stringlength(1|255);"))
		})

		It("should return status code of 422 if invalid id", func() {
			id := "abc123!@#xyz456"
			name := "Game Awesome Name"
			bundleID := "com.tfg.example"
			gameReader := JSONFor(JSON{
				"ID":       id,
				"Name":     name,
				"BundleID": bundleID,
			})
			request, _ := http.NewRequest("PUT", "/games", gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(BeEquivalentTo("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(ContainSubstring("ID: abc123!@#xyz456 does not validate as matches(^[^-][a-z0-9-]*$);"))
		})

		It("should return status code of 422 if empty id", func() {
			id := ""
			name := "Game Awesome Name"
			bundleID := "com.tfg.example"
			gameReader := JSONFor(JSON{
				"ID":       id,
				"Name":     name,
				"BundleID": bundleID,
			})
			request, _ := http.NewRequest("PUT", "/games", gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))

			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(BeEquivalentTo("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(ContainSubstring("ID: non zero value required;"))
		})

		It("should return status code of 500 if some error occurred", func() {
			gameReader := JSONFor(JSON{
				"ID":       uuid.NewV4().String(),
				"Name":     "Game Awesome Break",
				"BundleID": "com.topfreegames.example",
			})

			oldDB := app.DB
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			request, _ := http.NewRequest("PUT", "/games", gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("Upserting game failed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})
	})
})
