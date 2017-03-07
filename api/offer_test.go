// offers api
// https://github.com/topfreegames/offers
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
	"strings"

	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/topfreegames/offers/models"
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
			url := "/invalid"
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
		})
	})

	Describe("GET /offers", func() {
		It("should return available offers", func() {
			playerID := "player-1"
			gameID := "offers-game"
			url := fmt.Sprintf("/offers?player-id=%s&game-id=%s", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)
			var jsonBody map[string][]map[string]interface{}

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			err := json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(jsonBody).To(HaveKey("popup"))
			Expect(jsonBody).To(HaveKey("store"))
			popup := jsonBody["popup"]
			Expect(popup).To(HaveLen(1))
			Expect(popup[0]).To(HaveKey("id"))
			Expect(popup[0]).To(HaveKey("productId"))
			Expect(popup[0]).To(HaveKey("contents"))
			Expect(popup[0]).To(HaveKey("metadata"))
			store := jsonBody["store"]
			Expect(store).To(HaveLen(2))
			Expect(store[0]).To(HaveKey("id"))
			Expect(store[0]).To(HaveKey("productId"))
			Expect(store[0]).To(HaveKey("contents"))
			Expect(store[0]).To(HaveKey("metadata"))
			Expect(store[0]).To(HaveKey("expireAt"))
			Expect(store[1]).To(HaveKey("id"))
			Expect(store[1]).To(HaveKey("productId"))
			Expect(store[1]).To(HaveKey("contents"))
			Expect(store[1]).To(HaveKey("metadata"))
			Expect(store[1]).To(HaveKey("remainingPurchases"))
			Expect(store[1]).To(HaveKey("remainingImpressions"))
			Expect(store[1]).To(HaveKey("expireAt"))
		})

		It("should return empty list of available offers", func() {
			playerID := "player-1"
			gameID := "non-existing-offers-game"
			url := fmt.Sprintf("/offers?player-id=%s&game-id=%s", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)
			var jsonBody map[string]map[string]interface{}

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			err := json.Unmarshal(recorder.Body.Bytes(), &jsonBody)
			Expect(err).NotTo(HaveOccurred())
			Expect(jsonBody).To(BeEmpty())
			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

		It("should return status code 400 if player-id is not informed available offers", func() {
			gameID := "offers-game"
			url := fmt.Sprintf("/offers?game-id=%s", gameID)
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("The player-id parameter cannot be empty."))
			Expect(obj["description"]).To(Equal("The player-id parameter cannot be empty"))

		})

		It("should return status code 400 if game-id is not informed available offers", func() {
			playerID := "player-1"
			url := fmt.Sprintf("/offers?player-id=%s", playerID)
			request, _ := http.NewRequest("GET", url, nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
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
			url := fmt.Sprintf("/offers?player-id=%s&game-id=%s", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)

			oldDB := app.DB
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("Failed to retrieve offer for player"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})

		It("should not return offer after claim if offer template period has max 1", func() {
			// Create Offer by requesting it
			gameID := "limited-offers-game"
			playerID := "player-1"
			place := "store"
			url := fmt.Sprintf("/offers?player-id=%s&game-id=%s", playerID, gameID)
			request, _ := http.NewRequest("GET", url, nil)
			var body map[string][]*models.OfferToReturn

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))
			err := json.Unmarshal(recorder.Body.Bytes(), &body)
			Expect(err).ToNot(HaveOccurred())

			// Claim the Offer
			id := body[place][0].ID
			offerReader := JSONFor(JSON{
				"playerId": playerID,
				"gameId":   gameID,
			})
			request, _ = http.NewRequest("PUT", fmt.Sprintf("/offers/%s/claim", id), offerReader)
			recorder = httptest.NewRecorder()

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))

			// Offer must not be returned again in next Get
			request, _ = http.NewRequest("GET", url, nil)
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var newBody map[string][]*models.OfferToReturn
			err = json.Unmarshal(recorder.Body.Bytes(), &newBody)
			Expect(err).ToNot(HaveOccurred())
			Expect(newBody).NotTo(HaveKey(place))
		})
	})

	Describe("PUT /offers/{id}/claim", func() {
		It("should claim valid offer", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			gameID := "offers-game"
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   gameID,
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/claim", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Body.String()).To(Equal(fmt.Sprintf(`{"contents":{"gems":5,"gold":100},"nextAt":%v}`, app.Clock.GetTime().Unix()+1)))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			offer, err := models.GetOfferByID(app.DB, gameID, id, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offer).NotTo(BeNil())
			Expect(offer.ClaimedAt.Time.Unix()).To(Equal(app.Clock.GetTime().Unix()))
			Expect(offer.BoughtCounter).To(Equal(1))
		})

		It("should not claim a claimed offer", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			gameID := "offers-game"
			json := JSON{
				"playerId": "player-1",
				"gameId":   gameID,
			}
			offerReader1 := JSONFor(json)
			offerReader2 := JSONFor(json)
			request1, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/claim", id), offerReader1)
			request2, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/claim", id), offerReader2)

			app.Router.ServeHTTP(recorder, request1)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request2)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Body.String()).To(Equal(`{"contents":{"gems":5,"gold":100}}`))
			Expect(recorder.Code).To(Equal(http.StatusConflict))

			offer, err := models.GetOfferByID(app.DB, gameID, id, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offer).NotTo(BeNil())
			Expect(offer.ClaimedAt.Time.Unix()).To(Equal(app.Clock.GetTime().Unix()))
			Expect(offer.BoughtCounter).To(Equal(1))
		})

		It("should return 422 if invalid OfferID", func() {
			id := "invalid-offer-id"
			offerReader := JSONFor(JSON{
				"gameId":   "offers-game",
				"playerId": "player-1",
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/claim", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal(fmt.Sprintf("ID: %s does not validate;", id)))
		})

		It("should return 422 if missing parameters", func() {
			offerReader := JSONFor(JSON{})
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/claim", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("GameID: non zero value required;PlayerID: non zero value required;"))
		})

		It("should return 404 if non existing OfferID", func() {
			id := uuid.NewV4().String()
			offerReader := JSONFor(JSON{
				"gameId":   "offers-game",
				"playerId": "player-1",
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/claim", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 404 if non existing GameID", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerReader := JSONFor(JSON{
				"gameId":   "non-existing-offers-game",
				"playerId": "player-1",
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/claim", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferNotFoundError"))
			Expect(obj["description"]).To(Equal("Offer was not found with specified filters."))
		})

		It("should return 404 if non existing PlayerID", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerReader := JSONFor(JSON{
				"gameId":   "offers-game",
				"playerId": "non-existing-player-1",
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/claim", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferNotFoundError"))
			Expect(obj["description"]).To(Equal("Offer was not found with specified filters."))
		})

		It("should return status code of 500 if some error occurred", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerReader := JSONFor(JSON{
				"gameId":   "offers-game",
				"playerId": "player-1",
			})

			oldDB := app.DB
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/claim", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

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
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			invalidJSON := `"gameId:   "offers-g
											"Player-1"`
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/claim", id), strings.NewReader(invalidJSON))

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("json: cannot unmarshal string into Go value of type models.OfferToUpdate"))
		})

		It("should not claim id passed on body, only from url", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			bodyID := "35df52e7-3161-446f-975b-92f32871e37c"
			gameID := "offers-game"
			offerReader := JSONFor(JSON{
				"id":       bodyID,
				"playerId": "player-1",
				"gameId":   gameID,
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/offers/%s/claim", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			//Check if id was claimed and bodyID not
			claimedOffer, err := models.GetOfferByID(app.DB, gameID, id, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(claimedOffer.ClaimedAt.Valid).To(BeTrue())

			NotClaimedOffer, err := models.GetOfferByID(app.DB, gameID, bodyID, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(NotClaimedOffer.ClaimedAt.Valid).To(BeFalse())
		})

		It("should return status code 301 if empty id", func() {
			gameID := "offers-game"
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   gameID,
			})
			request, _ := http.NewRequest("PUT", "/offers//claim", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusMovedPermanently))
		})
	})

	Describe("POST /offers/{id}/impressions", func() {
		It("should update last seen at of valid offer", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			gameID := "offers-game"
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   gameID,
			})
			request, _ := http.NewRequest("POST", fmt.Sprintf("/offers/%s/impressions", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(int64(obj["nextAt"].(float64))).To(Equal(app.Clock.GetTime().Unix() + 1))
			offer, err := models.GetOfferByID(app.DB, gameID, id, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offer).NotTo(BeNil())
			Expect(offer.LastSeenAt.Time.Unix()).To(Equal(app.Clock.GetTime().Unix()))
			Expect(offer.SeenCounter).To(Equal(1))
		})

		It("should return the current timestamp as nextAt if offer reached max period", func() {
			id := "5ba8848f-1df0-45b3-b8b1-27a7d5eedd6a"
			gameID := "limited-offers-game"
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   gameID,
			})
			request, _ := http.NewRequest("POST", fmt.Sprintf("/offers/%s/impressions", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj).To(HaveLen(0))
			offer, err := models.GetOfferByID(app.DB, gameID, id, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offer).NotTo(BeNil())
			Expect(offer.LastSeenAt.Time.Unix()).To(Equal(app.Clock.GetTime().Unix()))
			Expect(offer.SeenCounter).To(Equal(1))
		})

		It("should return nextAt equal to now if offer has no every in frequency", func() {
			id := "4407b770-5b24-4ffa-8563-0694d1a10156"
			gameID := "offers-game"
			offerReader := JSONFor(JSON{
				"playerId": "player-11",
				"gameId":   gameID,
			})
			request, _ := http.NewRequest("POST", fmt.Sprintf("/offers/%s/impressions", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(int64(obj["nextAt"].(float64))).To(Equal(app.Clock.GetTime().Unix()))
			offer, err := models.GetOfferByID(app.DB, gameID, id, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(offer).NotTo(BeNil())
			Expect(offer.LastSeenAt.Time.Unix()).To(Equal(app.Clock.GetTime().Unix()))
			Expect(offer.SeenCounter).To(Equal(1))
		})

		It("should return status code 422 if invalid parameters", func() {
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   "offers-game",
			})
			url := "/offers/invalid-uuid/impressions"
			request, _ := http.NewRequest("POST", url, offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: invalid-uuid does not validate;"))
		})

		It("should return status code 422 if missing body parameters", func() {
			offerReader := JSONFor(JSON{})
			id := uuid.NewV4().String()
			request, _ := http.NewRequest("POST", fmt.Sprintf("/offers/%s/impressions", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("GameID: non zero value required;PlayerID: non zero value required;"))
		})

		It("should return status code 404 if offer with given ID does not exist", func() {
			id := uuid.NewV4().String()
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   "offers-game",
			})
			request, _ := http.NewRequest("POST", fmt.Sprintf("/offers/%s/impressions", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferNotFoundError"))
			Expect(obj["description"]).To(Equal("Offer was not found with specified filters."))
		})

		It("should return status code of 500 if some error occurred", func() {
			id := "56fc0477-39f1-485c-898e-4909e9155eb1"
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   "offers-game",
			})

			oldDB := app.DB
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			request, _ := http.NewRequest("POST", fmt.Sprintf("/offers/%s/impressions", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("sql: database is closed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})

		It("should return status code 404 if offer with given ID does not exist", func() {
			id := uuid.NewV4().String()
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   "offers-game",
			})
			request, _ := http.NewRequest("POST", fmt.Sprintf("/offers/%s/impressions", id), offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferNotFoundError"))
			Expect(obj["description"]).To(Equal("Offer was not found with specified filters."))
		})

		It("should return status code 301 if empty id", func() {
			gameID := "offers-game"
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   gameID,
			})
			request, _ := http.NewRequest("PUT", "/offers//impressions", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusMovedPermanently))
		})
	})
})
