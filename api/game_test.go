// offers api
// https://github.com/topfree/ames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	"encoding/json"
	"fmt"
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

	Describe("PUT /games/{id}", func() {
		It("should return status code of 200", func() {
			id := uuid.NewV4().String()
			gameReader := JSONFor(JSON{
				"Name":     "Game Awesome Name",
				"BundleID": "com.topfreegames.example",
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/games/%s", id), gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["gameId"]).To(Equal(id))
		})

		It("should return status code of 422 if missing parameter", func() {
			id := uuid.NewV4().String()
			gameReader := JSONFor(JSON{
				"Name": "Game Awesome Name",
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/games/%s", id), gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("BundleID: non zero value required;"))
		})

		It("should return status code of 422 if invalid name", func() {
			reallyBigName := "1234567890"
			for i := 0; i < 5; i++ {
				reallyBigName += reallyBigName
			}

			id := "game-id"
			gameReader := JSONFor(JSON{
				"Name":     reallyBigName,
				"BundleID": "com.topfreegames.example",
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/games/%s", id), gameReader)

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

			id := "game-id"
			gameReader := JSONFor(JSON{
				"Name":     uuid.NewV4().String(),
				"BundleID": reallyBigName,
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/games/%s", id), gameReader)

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
			id := "abc123!@$xyz456"
			name := "Game Awesome Name"
			bundleID := "com.tfg.example"
			gameReader := JSONFor(JSON{
				"Name":     name,
				"BundleID": bundleID,
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/games/%s", id), gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(BeEquivalentTo("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(ContainSubstring("ID: abc123!@$xyz456 does not validate;"))
		})

		It("should return status code of 422 if invalid id, even if is passed in body", func() {
			id := "abc123!@$xyz456"
			name := "Game Awesome Name"
			bundleID := "com.tfg.example"
			gameReader := JSONFor(JSON{
				"ID":       id,
				"Name":     name,
				"BundleID": bundleID,
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/games/%s", id), gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(BeEquivalentTo("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(ContainSubstring("ID: abc123!@$xyz456 does not validate;"))
		})

		It("should return status code of 404 if id is not passed", func() {
			gameReader := JSONFor(JSON{
				"Name":     "Game Awesome Name",
				"BundleID": "com.topfreegames.example",
			})
			request, _ := http.NewRequest("PUT", "/games/", gameReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			Expect(recorder.Body.String()).To(Equal("404 page not found\n"))
		})

		It("should return status code of 500 if some error occurred", func() {
			id := uuid.NewV4().String()
			gameReader := JSONFor(JSON{
				"Name":     "Game Awesome Break",
				"BundleID": "com.topfreegames.example",
			})

			oldDB := app.DB
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/games/%s", id), gameReader)

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

	Describe("GET /games", func() {
		It("should return status code of 200 and a list of games", func() {
			request, _ := http.NewRequest("GET", "/games", nil)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj []map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj).To(HaveLen(5))
			for i := 0; i < len(obj); i++ {
				Expect(obj[i]).To(HaveKey("id"))
				Expect(obj[i]).To(HaveKey("name"))
				Expect(obj[i]).To(HaveKey("bundleId"))
				Expect(obj[i]).To(HaveKey("metadata"))
			}
		})

		It("should return empty list if no games", func() {
			_, err := app.DB.DeleteFrom("offers").Exec()
			Expect(err).NotTo(HaveOccurred())
			_, err = app.DB.DeleteFrom("offer_templates").Exec()
			Expect(err).NotTo(HaveOccurred())
			_, err = app.DB.DeleteFrom("games").Exec()
			Expect(err).NotTo(HaveOccurred())
			request, _ := http.NewRequest("GET", "/games", nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(recorder.Body.String()).To(Equal("[]"))
		})

		It("should return status code of 500 if some error occurred", func() {
			oldDB := app.DB
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			request, _ := http.NewRequest("GET", "/games", nil)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("List games failed."))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})
	})
})
