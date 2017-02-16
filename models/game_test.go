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
			gameID := "game-id"

			var game models.Game
			err := db.
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
			meta := dat.JSON([]byte(`{"qwe": 123}`))
			game := &models.Game{
				Name:     "Game Awesome Name",
				ID:       "game-id-2",
				BundleID: "com.topfreegames.example",
				Metadata: meta,
			}
			err := db.
				InsertInto("games").
				Columns("name", "id", "bundle_id", "metadata").
				Record(game).
				Returning("created_at", "updated_at").
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
			gameID := "game-id"
			game, err := models.GetGameByID(db, gameID, nil)

			Expect(err).NotTo(HaveOccurred())
			Expect(game.ID).To(Equal(gameID))
			Expect(game.Name).To(Equal("game-1"))

			obj, err := game.GetMetadata()
			Expect(err).NotTo(HaveOccurred())
			Expect(obj).To(BeEquivalentTo(map[string]interface{}{}))
		})

		It("Should return error if game not found", func() {
			gameID := uuid.NewV4().String()
			expectedError := errors.NewModelNotFoundError("Game", map[string]interface{}{
				"ID": gameID,
			})
			game, err := models.GetGameByID(db, gameID, nil)
			Expect(game).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expectedError))
		})
	})

	Describe("Upsert game", func() {
		It("should insert game with new id", func() {
			id := uuid.NewV4().String()
			game := models.Game{
				ID:       id,
				Name:     "Game Awesome Name",
				BundleID: "com.tfg.example",
			}
			err := models.UpsertGame(db, &game, nil)
			Expect(err).NotTo(HaveOccurred())

			gameFromDB, err := models.GetGameByID(db, id, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(gameFromDB.ID).To(Equal(id))
		})

		It("should update game with existing id", func() {
			id := "upsert-game-id"
			name := "Game Awesome Name"
			bundleID := "com.tfg.example"
			game := models.Game{
				ID:       id,
				Name:     name,
				BundleID: bundleID,
			}
			err := models.UpsertGame(db, &game, nil)
			Expect(err).NotTo(HaveOccurred())

			gameFromDB, err := models.GetGameByID(db, id, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(gameFromDB.ID).To(Equal(id))
			Expect(gameFromDB.Name).To(Equal(name))
			Expect(gameFromDB.BundleID).To(Equal(bundleID))
		})
	})
})
