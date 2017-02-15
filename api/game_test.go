package api

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topfreegames/offers/api"
	"testing"
)

var _ = Describe("Game Handler", func() {
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		// Record HTTP responses.
		recorder = httptest.NewRecorder()
	})

	Describe("POST /game", func() {
		It("should return status code of 200", func() {
			game := models.Game{
				Name:     uuid.NewV4().String(),
				Metagame: nil,
			}

			b, _ := json.Marshal(&game)

			request, _ = http.NewRequest("POST", "/game/upsert", b)
			app.Router.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(200))
		})
	})
})
