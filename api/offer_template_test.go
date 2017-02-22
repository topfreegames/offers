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

	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	. "github.com/topfreegames/offers/testing"
)

var _ = Describe("Offer Template Handler", func() {
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	Describe("PUT /offer-templates", func() {
		It("should return status code 200 for valid parameters", func() {
			id := uuid.NewV4().String()
			offerTemplateReader := JSONFor(JSON{
				"ID":        id,
				"Name":      "New Awesome Game",
				"ProductID": "com.tfg.example",
				"GameID":    "game-id",
				"Contents":  dat.JSON([]byte("{\"gems\": 5, \"gold\": 100}")),
				"Period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"Frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"Trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"Placement": "popup",
			})

			request, _ := http.NewRequest("PUT", "/offer-templates", offerTemplateReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

		It("should return status code 422 if ID not passed", func() {
			offerTemplateReader := JSONFor(JSON{
				"Name":      "New Awesome Game",
				"ProductID": "com.tfg.example",
				"GameID":    "game-id",
				"Contents":  dat.JSON([]byte("{\"gems\": 5, \"gold\": 100}")),
				"Period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"Frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"Trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"Placement": "popup",
			})

			request, _ := http.NewRequest("PUT", "/offer-templates", offerTemplateReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
		})

		It("should return status code 422 if ID is empty", func() {
			offerTemplateReader := JSONFor(JSON{
				"ID":        "",
				"Name":      "New Awesome Game",
				"ProductID": "com.tfg.example",
				"GameID":    "game-id",
				"Contents":  dat.JSON([]byte("{\"gems\": 5, \"gold\": 100}")),
				"Period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"Frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"Trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"Placement": "popup",
			})

			request, _ := http.NewRequest("PUT", "/offer-templates", offerTemplateReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
		})

		It("should return status code 422 if ID is invalid", func() {
			offerTemplateReader := JSONFor(JSON{
				"ID":        "not_a_valid id!",
				"Name":      "New Awesome Game",
				"ProductID": "com.tfg.example",
				"GameID":    "game-id",
				"Contents":  dat.JSON([]byte("{\"gems\": 5, \"gold\": 100}")),
				"Period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"Frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"Trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"Placement": "popup",
			})

			request, _ := http.NewRequest("PUT", "/offer-templates", offerTemplateReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
		})

		It("should return status code 422 if game-id doesn`t exist", func() {
			id := uuid.NewV4().String()
			offerTemplateReader := JSONFor(JSON{
				"ID":        id,
				"ProductID": "com.tfg.example",
				"GameID":    "not-existing-game-id",
				"Contents":  dat.JSON([]byte("{\"gems\": 5, \"gold\": 100}")),
				"Period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"Frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"Trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"Placement": "popup",
			})

			request, _ := http.NewRequest("PUT", "/offer-templates", offerTemplateReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
		})

		It("returns status code of 500 if database is unavailable", func() {
			id := uuid.NewV4().String()
			offerTemplateReader := JSONFor(JSON{
				"ID":        id,
				"Name":      "New Awesome Game",
				"ProductID": "com.tfg.example",
				"GameID":    "game-id",
				"Contents":  dat.JSON([]byte("{\"gems\": 5, \"gold\": 100}")),
				"Period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"Frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"Trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"Placement": "popup",
			})

			request, _ := http.NewRequest("PUT", "/offer-templates", offerTemplateReader)

			oldDB := app.DB
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			Expect(recorder.Body.String()).To(Equal("Insert offer template failed"))
			app.DB = oldDB // avoid errors in after each
		})
	})
})
