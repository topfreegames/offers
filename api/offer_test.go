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
	"strings"

	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	. "github.com/topfreegames/offers/testing"
)

var _ = Describe("Offer Handler", func() {
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		// Record HTTP responses.
		recorder = httptest.NewRecorder()
	})

	Describe("Invalid route", func() {
		It("should return status code 404 if invalid route", func() {
			//Given
			url := "/invalid"
			request, _ := http.NewRequest("GET", url, nil)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("GET /offers", func() {
		It("should return available offers", func() {
			//Given
			playerID := "player-1"
			gameID := "offers-game"
			url := "/offers?player-id=" + playerID + "&game-id=" + gameID
			request, _ := http.NewRequest("GET", url, nil)
			var jsonBody map[string][]map[string]interface{}

			//When
			app.Router.ServeHTTP(recorder, request)
			err := json.Unmarshal(recorder.Body.Bytes(), &jsonBody)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonBody).To(HaveKey("popup"))
			Expect(jsonBody).To(HaveKey("store"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

		It("should return empty list of available offers", func() {
			//Given
			playerID := "player-1"
			gameID := "non-existing-offers-game"
			url := "/offers?player-id=" + playerID + "&game-id=" + gameID
			request, _ := http.NewRequest("GET", url, nil)
			var jsonBody map[string]map[string]interface{}

			//When
			app.Router.ServeHTTP(recorder, request)
			err := json.Unmarshal(recorder.Body.Bytes(), &jsonBody)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonBody).To(BeEmpty())
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

		It("should return status code 400 if player-id is not informed available offers", func() {
			//Given
			gameID := "offers-game"
			url := "/offers?game-id=" + gameID
			request, _ := http.NewRequest("GET", url, nil)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("The player-id parameter cannot be empty."))
			Expect(obj["description"]).To(Equal("The player-id parameter cannot be empty"))

		})

		It("should return status code 400 if game-id is not informed available offers", func() {
			//Given
			playerID := "player-1"
			url := "/offers?player-id=" + playerID
			request, _ := http.NewRequest("GET", url, nil)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("The game-id parameter cannot be empty."))
			Expect(obj["description"]).To(Equal("The game-id parameter cannot be empty"))
		})

		It("should return status code of 500 if some error occurred", func() {
			playerID := "player-1"
			gameID := "offers-game"
			url := "/offers?player-id=" + playerID + "&game-id=" + gameID
			request, _ := http.NewRequest("GET", url, nil)

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
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("Failed to retrieve offer for player"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})
	})

	Describe("PUT /offer/claim", func() {
		It("should claim valid offer", func() {
			//Given
			offerReader := JSONFor(JSON{
				"ID":       "56fc0477-39f1-485c-898e-4909e9155eb1",
				"GameID":   "offers-game",
				"PlayerID": "player-1",
			})
			request, _ := http.NewRequest("PUT", "/offer/claim", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Body.String()).To(Equal(`{"gems": 5, "gold": 100}`))
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

		It("should not claim a claimed offer", func() {
			//Given
			offerReader1 := JSONFor(JSON{
				"ID":       "56fc0477-39f1-485c-898e-4909e9155eb1",
				"GameID":   "offers-game",
				"PlayerID": "player-1",
			})
			offerReader2 := JSONFor(JSON{
				"ID":       "56fc0477-39f1-485c-898e-4909e9155eb1",
				"GameID":   "offers-game",
				"PlayerID": "player-1",
			})
			request1, _ := http.NewRequest("PUT", "/offer/claim", offerReader1)
			request2, _ := http.NewRequest("PUT", "/offer/claim", offerReader2)

			//When
			app.Router.ServeHTTP(recorder, request1)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request2)

			//Then
			Expect(recorder.Body.String()).To(Equal(`{"gems": 5, "gold": 100}`))
			Expect(recorder.Code).To(Equal(http.StatusConflict))
		})

		It("should return 422 if invalid OfferID", func() {
			//Given
			offerReader := JSONFor(JSON{
				"ID":       "567-391-4c-8-4909eeb1",
				"GameID":   "offers-game",
				"PlayerID": "player-1",
			})
			request, _ := http.NewRequest("PUT", "/offer/claim", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Body.String()).To(Equal(`{"code":"OFF-002","description":"ID: 567-391-4c-8-4909eeb1 does not validate as uuidv4;","error":"ValidationFailedError"}`))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: 567-391-4c-8-4909eeb1 does not validate as uuidv4;"))
		})

		It("should return 404 if non existing OfferID", func() {
			//Given
			offerReader := JSONFor(JSON{
				"ID":       uuid.NewV4().String(),
				"GameID":   "offers-game",
				"PlayerID": "player-1",
			})
			request, _ := http.NewRequest("PUT", "/offer/claim", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 404 if non existing GameID", func() {
			//Given
			offerReader := JSONFor(JSON{
				"ID":       "56fc0477-39f1-485c-898e-4909e9155eb1",
				"GameID":   "non-existing-offers-game",
				"PlayerID": "player-1",
			})
			request, _ := http.NewRequest("PUT", "/offer/claim", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferNotFoundError"))
			Expect(obj["description"]).To(Equal("Offer was not found with specified filters."))
		})

		It("should return 404 if non existing PlayerID", func() {
			//Given
			offerReader := JSONFor(JSON{
				"ID":       "56fc0477-39f1-485c-898e-4909e9155eb1",
				"GameID":   "offers-game",
				"PlayerID": "non-existing-player-1",
			})
			request, _ := http.NewRequest("PUT", "/offer/claim", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferNotFoundError"))
			Expect(obj["description"]).To(Equal("Offer was not found with specified filters."))
		})

		It("should return 422 if OfferID is not passed", func() {
			//Given
			offerReader := JSONFor(JSON{
				"GameID":   "offers-game",
				"PlayerID": "non-existing-player-1",
			})
			request, _ := http.NewRequest("PUT", "/offer/claim", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: non zero value required;"))
		})

		It("should return status code of 500 if some error occurred", func() {
			offerReader := JSONFor(JSON{
				"id":       "56fc0477-39f1-485c-898e-4909e9155eb1",
				"gameId":   "offers-game",
				"playerId": "player-1",
			})

			oldDB := app.DB
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			request, _ := http.NewRequest("PUT", "/offer/claim", offerReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("sql: database is closed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})

		It("should return status code 400 if invalid json is sent", func() {
			//Given
			invalidJSON := `"ID    "56fc047-39f1-485c-898e-4909e9155eb1",
											"GameID:   "offers-g
											"Player-1"`
			request, _ := http.NewRequest("PUT", "/offer/claim", strings.NewReader(invalidJSON))

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("json: cannot unmarshal string into Go value of type models.OfferToUpdate"))
		})
	})

	Describe("PUT /update-offer-last-seen-at", func() {
		It("should update last seen at of valid offer", func() {
			//Given
			offerReader := JSONFor(JSON{
				"ID":       "56fc0477-39f1-485c-898e-4909e9155eb1",
				"PlayerID": "player-1",
				"GameID":   "offers-game",
			})
			request, _ := http.NewRequest("PUT", "/offer/last-seen-at", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

		It("should return status code 422 if ID is invalid", func() {
			//Given
			offerReader := JSONFor(JSON{
				"ID":       "invalid-uuid",
				"PlayerID": "player-1",
				"GameID":   "offers-game",
			})
			request, _ := http.NewRequest("PUT", "/offer/last-seen-at", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: invalid-uuid does not validate as uuidv4;"))
		})

		It("should return status code 422 if ID is not passed", func() {
			//Given
			offerReader := JSONFor(JSON{
				"PlayerID": "player-1",
				"GameID":   "offers-game",
			})
			request, _ := http.NewRequest("PUT", "/offer/last-seen-at", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Body.String()).To(Equal(`{"code":"OFF-002","description":"ID: non zero value required;","error":"ValidationFailedError"}`))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: non zero value required;"))
		})

		It("should return status code 404 if offer with given ID does not exist", func() {
			//Given
			offerReader := JSONFor(JSON{
				"ID":       uuid.NewV4().String(),
				"PlayerID": "player-1",
				"GameID":   "offers-game",
			})
			request, _ := http.NewRequest("PUT", "/offer/last-seen-at", offerReader)

			//When
			app.Router.ServeHTTP(recorder, request)

			//Then
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferNotFoundError"))
			Expect(obj["description"]).To(Equal("Offer was not found with specified filters."))
		})

		It("should return status code of 500 if some error occurred", func() {
			offerReader := JSONFor(JSON{
				"ID":       "56fc0477-39f1-485c-898e-4909e9155eb1",
				"PlayerID": "player-1",
				"GameID":   "offers-game",
			})

			oldDB := app.DB
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			request, _ := http.NewRequest("PUT", "/offer/last-seen-at", offerReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("sql: database is closed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})
	})
})
