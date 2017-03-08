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

	"gopkg.in/mgutz/dat.v2/dat"
	runner "gopkg.in/mgutz/dat.v2/sqlx-runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/topfreegames/offers/models"
	. "github.com/topfreegames/offers/testing"
)

var _ = Describe("Offer Template Handler", func() {
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
	})

	Describe("POST /templates/{key}/insert", func() {
		It("should return status code 201 for valid parameters", func() {
			name := "New Awesome Game"
			key := uuid.NewV4().String()
			productID := "com.tfg.example"
			gameID := "game-id"
			contents := "{\"gems\": 5, \"gold\": 100}"
			period := "{\"max\": 1}"
			frequency := "{\"every\": \"24h\"}"
			trigger := "{\"from\": 1487280506875, \"to\": 1487366964730}"
			placement := "popup"
			offerTemplateReader := JSONFor(JSON{
				"name":      name,
				"key":       key,
				"productId": productID,
				"gameId":    gameID,
				"contents":  dat.JSON([]byte(contents)),
				"period":    dat.JSON([]byte(period)),
				"frequency": dat.JSON([]byte(frequency)),
				"trigger":   dat.JSON([]byte(trigger)),
				"placement": placement,
			})

			request, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/insert", key), offerTemplateReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusCreated), recorder.Body.String())
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).NotTo(BeEmpty())
			Expect(obj["key"]).NotTo(BeEmpty())
			Expect(obj["name"]).To(Equal(name))
			Expect(obj["productId"]).To(Equal(productID))
			Expect(obj["gameId"]).To(Equal(gameID))
			Expect(obj["contents"].(map[string]interface{})["gems"].(float64)).To(BeEquivalentTo(5))
			Expect(obj["contents"].(map[string]interface{})["gold"]).To(BeEquivalentTo(100))
			Expect(obj["period"].(map[string]interface{})["max"].(float64)).To(BeEquivalentTo(1))
			Expect(obj["frequency"].(map[string]interface{})["every"].(string)).To(Equal("24h"))
			Expect(obj["trigger"].(map[string]interface{})["to"].(float64)).To(BeEquivalentTo(1487366964730))
			Expect(obj["trigger"].(map[string]interface{})["from"].(float64)).To(BeEquivalentTo(1487280506875))
			Expect(obj["placement"]).To(Equal(placement))
			Expect(obj["enabled"]).To(BeTrue())
		})

		It("should return status code 422 if missing arguments", func() {
			offerTemplateReader := JSONFor(JSON{})

			request, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/insert", uuid.NewV4().String()), offerTemplateReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("Name: non zero value required;ProductID: non zero value required;GameID: non zero value required;Contents: [] does not validate as RequiredJSONObject;;Period: [] does not validate as RequiredJSONObject;;Frequency: [] does not validate as RequiredJSONObject;;Trigger: [] does not validate as RequiredJSONObject;;Placement: non zero value required;"))
		})

		It("should return status code 422 if invalid arguments", func() {
			offerTemplateReader := JSONFor(JSON{
				"name":      "",
				"productId": "",
				"gameId":    "___",
				"contents":  "{not-a-json}",
				"period":    "{not-a-json}",
				"frequency": "{not-a-json}",
				"trigger":   "{not-a-json}",
				"placement": "",
			})

			request, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/insert", uuid.NewV4().String()), offerTemplateReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("Name: non zero value required;ProductID: non zero value required;GameID: ___ does not validate as matches(^[^-][a-z0-9-]*$);Contents: [34 123 110 111 116 45 97 45 106 115 111 110 125 34] does not validate as RequiredJSONObject;;Period: [34 123 110 111 116 45 97 45 106 115 111 110 125 34] does not validate as RequiredJSONObject;;Frequency: [34 123 110 111 116 45 97 45 106 115 111 110 125 34] does not validate as RequiredJSONObject;;Trigger: [34 123 110 111 116 45 97 45 106 115 111 110 125 34] does not validate as RequiredJSONObject;;Placement: non zero value required;"))
		})

		It("should return status code 422 if game-id doesn`t exist", func() {
			key := uuid.NewV4().String()
			offerTemplateReader := JSONFor(JSON{
				"name":      "New Awesome Game",
				"productId": "com.tfg.example",
				"gameId":    "not-existing-game-id",
				"contents":  dat.JSON([]byte("{\"gems\": 5, \"gold\": 100}")),
				"period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"placement": "popup",
			})

			request, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/insert", key), offerTemplateReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity), recorder.Body.String())
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-003"))
			Expect(obj["error"]).To(Equal("InvalidOfferTemplateError"))
			Expect(obj["description"]).To(Equal("OfferTemplate could not be saved due to: insert or update on table \"offer_templates\" violates foreign key constraint \"offer_templates_game_id_fkey\""))
		})

		It("should return status code 422 if contents is empty", func() {
			key := uuid.NewV4().String()
			offerTemplateReader := JSONFor(JSON{
				"name":      "New Awesome Game",
				"productId": "com.tfg.example",
				"gameId":    "game-id",
				"contents":  "",
				"period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"placement": "popup",
			})

			request, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/insert", key), offerTemplateReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity), recorder.Body.String())
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("Contents: [34 34] does not validate as RequiredJSONObject;;"))
		})

		It("should return status code 301 if key is not sent", func() {
			offerTemplateReader := JSONFor(JSON{
				"name":      "New Awesome Game",
				"productId": "com.tfg.example",
				"gameId":    "game-id",
				"contents":  dat.JSON([]byte("{\"gems\": 5, \"gold\": 100}")),
				"period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"placement": "popup",
			})

			request, _ := http.NewRequest("POST", "/templates//insert", offerTemplateReader)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusMovedPermanently), recorder.Body.String())
		})

		It("returns status code of 500 if database is unavailable", func() {
			key := uuid.NewV4().String()
			offerTemplateReader := JSONFor(JSON{
				"name":      "New Awesome Game",
				"productId": "com.tfg.example",
				"gameId":    "game-id",
				"contents":  dat.JSON([]byte("{\"gems\": 5, \"gold\": 100}")),
				"period":    dat.JSON([]byte("{\"type\": \"once\"}")),
				"frequency": dat.JSON([]byte("{\"every\": 24, \"unit\": \"hour\"}")),
				"trigger":   dat.JSON([]byte("{\"from\": 1487280506875, \"to\": 1487366964730}")),
				"placement": "popup",
			})

			request, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/insert", key), offerTemplateReader)

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
			Expect(obj["error"]).To(Equal("Insert offer template failed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})

		It("should return status code 201 if pair offer template key and game already exists and is disabled", func() {
			templateID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerTemplateReaderEnable := JSONFor(JSON{})
			name := "template-1"
			key := "da700673-0415-43c3-a8e0-18331b794482"
			productID := "com.tfg.example"
			gameID := "offers-game"
			contents := "{\"gems\": 5, \"gold\": 100}"
			period := "{\"max\": 1}"
			frequency := "{\"every\": \"24h\"}"
			trigger := "{\"from\": 1487280506875, \"to\": 1487366964730}"
			placement := "popup"
			offerTemplateReader := JSONFor(JSON{
				"name":      name,
				"productId": productID,
				"gameId":    gameID,
				"contents":  dat.JSON([]byte(contents)),
				"period":    dat.JSON([]byte(period)),
				"frequency": dat.JSON([]byte(frequency)),
				"trigger":   dat.JSON([]byte(trigger)),
				"placement": placement,
			})

			request1, _ := http.NewRequest("PUT", fmt.Sprintf("/templates/%s/disable", templateID), offerTemplateReaderEnable)
			request2, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/insert", key), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request1)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))

			recorder = httptest.NewRecorder()
			app.Router.ServeHTTP(recorder, request2)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusCreated))
		})

		It("should return status code 409 for invalid key", func() {
			key := "this-is-an-invalid-uuid"
			offerTemplateReader := JSONFor(JSON{
				"name":      "New Awesome Game",
				"productId": "com.tfg.example",
				"gameId":    "game-id",
				"contents":  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				"period":    dat.JSON([]byte(`{"max": 1}`)),
				"frequency": dat.JSON([]byte(`{"every": "24h"}`)),
				"trigger":   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
				"placement": "popup",
			})
			request, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/insert", key), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity), recorder.Body.String())
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: this-is-an-invalid-uuid does not validate;"))
		})

		It("should return 409 if offer template is enabled and trying to insert new one with same key", func() {
			offerTemplateKey := "f8f8c92a-fb88-4a7d-8a5b-b8de66d7b5ac"
			offerTemplateReader := JSONFor(JSON{
				"name":      "template-to-update",
				"key":       offerTemplateKey,
				"productId": "com.tfg.example",
				"gameId":    "offers-game-to-update",
				"contents":  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				"period":    dat.JSON([]byte(`{"max": 1}`)),
				"frequency": dat.JSON([]byte(`{"every": "24h"}`)),
				"trigger":   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
				"placement": "popup",
			})
			request, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/insert", offerTemplateKey), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusConflict))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-003"))
			Expect(obj["error"]).To(Equal("ConflictedOfferTemplateError"))
			Expect(obj["description"]).To(Equal("OfferTemplate could not be saved due to: There is another enabled offer template with key f8f8c92a-fb88-4a7d-8a5b-b8de66d7b5ac"))
		})
	})

	Describe("POST /templates/{key}/update", func() {
		It("should return 201 if offer template with key is enabled and a new version is created", func() {
			offerTemplateID := "3479fa18-6739-4ea8-a011-24ae40ce1c02"
			offerTemplateKey := "f8f8c92a-fb88-4a7d-8a5b-b8de66d7b5ac"
			offerTemplateReader := JSONFor(JSON{
				"name":      "template-to-update",
				"key":       offerTemplateKey,
				"productId": "com.tfg.example",
				"gameId":    "offers-game-to-update",
				"contents":  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				"period":    dat.JSON([]byte(`{"max": 1}`)),
				"frequency": dat.JSON([]byte(`{"every": "24h"}`)),
				"trigger":   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
				"placement": "popup",
			})
			request, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/update", offerTemplateKey), offerTemplateReader)
			var newOfferTemplate models.OfferTemplate

			app.Router.ServeHTTP(recorder, request)

			oldOt, err := models.GetOfferTemplateByID(app.DB, offerTemplateID, nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(recorder.Code).To(Equal(http.StatusCreated))
			Expect(oldOt.Enabled).To(BeFalse())

			json.Unmarshal(recorder.Body.Bytes(), &newOfferTemplate)
			Expect(newOfferTemplate.Enabled).To(BeTrue())
		})

		It("should insert offer template if trying to update inexistent key", func() {
			offerTemplateKey := uuid.NewV4().String()
			offerTemplateReader := JSONFor(JSON{
				"name":      "template-to-update",
				"key":       offerTemplateKey,
				"productId": "com.tfg.example",
				"gameId":    "offers-game-to-update",
				"metadata":  dat.JSON([]byte(`{"meta": "data"}`)),
				"contents":  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				"period":    dat.JSON([]byte(`{"max": 1}`)),
				"frequency": dat.JSON([]byte(`{"every": "24h"}`)),
				"trigger":   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
				"placement": "popup",
			})
			request, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/update", offerTemplateKey), offerTemplateReader)
			var newOfferTemplate models.OfferTemplate

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusCreated))
			json.Unmarshal(recorder.Body.Bytes(), &newOfferTemplate)
			Expect(newOfferTemplate.Enabled).To(BeTrue())
		})

		It("should return status code 422 if contents isn't a JSON", func() {
			offerTemplateKey := uuid.NewV4().String()
			offerTemplateReader := JSONFor(JSON{
				"name":      "template-to-update",
				"key":       offerTemplateKey,
				"productId": "com.tfg.example",
				"gameId":    "offers-game-to-update",
				"contents":  5,
				"period":    dat.JSON([]byte(`{"max": 1}`)),
				"frequency": dat.JSON([]byte(`{"every": "24h"}`)),
				"trigger":   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
				"placement": "popup",
			})
			request, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/update", offerTemplateKey), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
		})

		It("should return status code 422 if metadata isn't a JSON", func() {
			offerTemplateKey := uuid.NewV4().String()
			offerTemplateReader := JSONFor(JSON{
				"name":      "template-to-update",
				"key":       offerTemplateKey,
				"productId": "com.tfg.example",
				"gameId":    "offers-game-to-update",
				"metadata":  "5",
				"contents":  dat.JSON([]byte(`{"gems": 5, "gold": 100}`)),
				"period":    dat.JSON([]byte(`{"max": 1}`)),
				"frequency": dat.JSON([]byte(`{"every": "24h"}`)),
				"trigger":   dat.JSON([]byte(`{"from": 1487280506875, "to": 1487366964730}`)),
				"placement": "popup",
			})
			request, _ := http.NewRequest("POST", fmt.Sprintf("/templates/%s/update", offerTemplateKey), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
		})
	})

	Describe("PUT /templates/{id}/enable", func() {
		It("should enable an enabled offer template", func() {
			templateID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerTemplateReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/templates/%s/enable", templateID), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).To(Equal(templateID))
		})

		It("should enable a disabled offer template", func() {
			templateID := "27b0370f-bd61-4346-a10d-50ec052ae125"
			offerTemplateReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/templates/%s/enable", templateID), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).To(Equal(templateID))
		})

		It("returns status code of 500 if database is unavailable", func() {
			templateID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerTemplateReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/templates/%s/enable", templateID), offerTemplateReader)

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
			Expect(obj["error"]).To(Equal("Update offer template failed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})

		It("should return status code 404 if id doesn't exist", func() {
			templateID := uuid.NewV4().String()
			offerTemplateReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/templates/%s/enable", templateID), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferTemplateNotFoundError"))
			Expect(obj["description"]).To(Equal("OfferTemplate was not found with specified filters."))
		})

		It("should return status code 422 if invalid parameters", func() {
			templateID := "not-uuid"
			offerTemplateReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/templates/%s/enable", templateID), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: not-uuid does not validate;"))
		})

		It("should return status code 301 if empty id", func() {
			gameID := "offers-game"
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   gameID,
			})
			request, _ := http.NewRequest("PUT", "/templates//enable", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusMovedPermanently))
		})
	})

	Describe("PUT /templates/{id}/disable", func() {
		It("should disable an enabled offer template", func() {
			templateID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerTemplateReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/templates/%s/disable", templateID), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).To(Equal(templateID))
		})

		It("should disable a disabled offer template", func() {
			templateID := "27b0370f-bd61-4346-a10d-50ec052ae125"
			offerTemplateReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/templates/%s/disable", templateID), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).To(Equal(templateID))
		})

		It("should use ID from URI even if a valid one is passed in body", func() {
			templateID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerTemplateReader := JSONFor(JSON{
				"id": "aa65a3f2-7cf8-4d76-957f-0a23a1bbbd32",
			})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/templates/%s/disable", templateID), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["id"]).To(Equal(templateID))
		})

		It("returns status code of 500 if database is unavailable", func() {
			templateID := "dd21ec96-2890-4ba0-b8e2-40ea67196990"
			offerTemplateReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/templates/%s/disable", templateID), offerTemplateReader)

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
			Expect(obj["error"]).To(Equal("Update offer template failed"))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})

		It("should return status code 404 if id doesn't exist", func() {
			templateID := uuid.NewV4().String()
			offerTemplateReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/templates/%s/disable", templateID), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusNotFound))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-001"))
			Expect(obj["error"]).To(Equal("OfferTemplateNotFoundError"))
			Expect(obj["description"]).To(Equal("OfferTemplate was not found with specified filters."))
		})

		It("should return status code 422 if invalid parameters", func() {
			templateID := "not-uuid"
			offerTemplateReader := JSONFor(JSON{})
			request, _ := http.NewRequest("PUT", fmt.Sprintf("/templates/%s/disable", templateID), offerTemplateReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusUnprocessableEntity))
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-002"))
			Expect(obj["error"]).To(Equal("ValidationFailedError"))
			Expect(obj["description"]).To(Equal("ID: not-uuid does not validate;"))
		})

		It("should return status code 301 if empty id", func() {
			gameID := "offers-game"
			offerReader := JSONFor(JSON{
				"playerId": "player-1",
				"gameId":   gameID,
			})
			request, _ := http.NewRequest("PUT", "/templates//disable", offerReader)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(http.StatusMovedPermanently))
		})
	})

	Describe("GET /templates", func() {
		It("should return status code of 200 and a list of offer templates", func() {
			request, _ := http.NewRequest("GET", "/templates?game-id=offers-game", nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusOK))
			var obj []map[string]interface{}
			err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj).To(HaveLen(4))
			for i := 0; i < len(obj); i++ {
				Expect(obj[i]).To(HaveKey("id"))
				Expect(obj[i]).To(HaveKey("name"))
				Expect(obj[i]).To(HaveKey("productId"))
				Expect(obj[i]).To(HaveKey("gameId"))
				Expect(obj[i]).To(HaveKey("contents"))
				Expect(obj[i]).To(HaveKey("metadata"))
				Expect(obj[i]).To(HaveKey("enabled"))
				Expect(obj[i]).To(HaveKey("placement"))
				Expect(obj[i]).To(HaveKey("period"))
				Expect(obj[i]).To(HaveKey("frequency"))
				Expect(obj[i]).To(HaveKey("trigger"))
			}
		})

		It("should return empty list if no offer templates", func() {
			request, _ := http.NewRequest("GET", "/templates?game-id=unexistent-game", nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(recorder.Code).To(Equal(http.StatusOK))
			Expect(recorder.Body.String()).To(Equal("[]"))
		})

		It("should return status code of 400 if game-id is not provided", func() {
			request, _ := http.NewRequest("GET", "/templates", nil)

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
			oldDB := app.DB
			db, err := GetTestDB()
			Expect(err).NotTo(HaveOccurred())
			app.DB = db
			app.DB.(*runner.DB).DB.Close() // make DB connection unavailable
			request, _ := http.NewRequest("GET", "/templates?game-id=offers-game", nil)

			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Header().Get("Content-Type")).To(Equal("application/json"))

			Expect(recorder.Code).To(Equal(http.StatusInternalServerError))
			var obj map[string]interface{}
			err = json.Unmarshal([]byte(recorder.Body.String()), &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["code"]).To(Equal("OFF-004"))
			Expect(obj["error"]).To(Equal("List game offer templates failed."))
			Expect(obj["description"]).To(Equal("sql: database is closed"))
			app.DB = oldDB // avoid errors in after each
		})
	})
})
