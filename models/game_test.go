// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	"github.com/mgutz/dat"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/topfreegames/offers/errors"
	"github.com/topfreegames/offers/models"
)

var _ = Describe("Games Model", func() {
	Describe("Game Instance", func() {
		It("Shoud load a game", func() {
			gameID, err := uuid.FromString("3393cd15-5b5a-4cfb-9725-ddbde660a727")
			Expect(err).NotTo(HaveOccurred())

			var game models.Game
			err = db.
				Select("*").
				From("games").
				Where("id = $1", gameID).
				QueryStruct(&game)

			Expect(err).NotTo(HaveOccurred())
			Expect(game.ID).To(Equal(gameID))
			Expect(game.Name).To(Equal("game-1"))

			obj, err := game.GetMetadata()
			Expect(err).NotTo(HaveOccurred())
			Expect(obj).To(BeEquivalentTo(map[string]interface{}{}))
		})

		It("Should create game", func() {
			meta := dat.JSON(
				[]byte(`
					{"qwe": 123}
				`),
			)
			id := uuid.NewV4().String()

			game := &models.Game{
				Name:     id,
				Metadata: meta,
				BundleID: "com.topfreegames.example",
			}
			err := db.
				InsertInto("games").
				Columns("name", "metadata", "bundle_id").
				Record(game).
				Returning("id", "created_at", "updated_at").
				QueryStruct(game)

			Expect(err).NotTo(HaveOccurred())
			Expect(game.ID).NotTo(Equal(""))

			var game2 models.Game
			err = db.
				Select("*").
				From("games").
				Where("id = $1", game.ID).
				QueryStruct(&game2)

			obj, err := game2.GetMetadata()
			Expect(err).NotTo(HaveOccurred())
			Expect(obj.(map[string]interface{})["qwe"]).To(BeEquivalentTo(123))
		})
	})

	Describe("Get game by id", func() {
		It("Should load game by id", func() {
			gameID, err := uuid.FromString("3393cd15-5b5a-4cfb-9725-ddbde660a727")
			Expect(err).NotTo(HaveOccurred())

			game, err := models.GetGameByID(db, gameID)
			Expect(err).NotTo(HaveOccurred())
			Expect(game.ID).To(Equal(gameID))
			Expect(game.Name).To(Equal("game-1"))

			obj, err := game.GetMetadata()
			Expect(err).NotTo(HaveOccurred())
			Expect(obj).To(BeEquivalentTo(map[string]interface{}{}))
		})

		It("Should return error if game not found", func() {
			gameID := uuid.NewV4()
			expectedError := errors.NewGameNotFoundError(map[string]interface{}{
				"ID": gameID,
			})
			game, err := models.GetGameByID(db, gameID)
			Expect(game).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expectedError))
		})
	})

	Describe("Upsert game", func() {
		It("Should upsert game", func() {
		})
	})
})
