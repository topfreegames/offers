// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/topfreegames/offers/errors"
	"github.com/topfreegames/offers/models"
	"gopkg.in/mgutz/dat.v2/dat"
)

var _ = Describe("Games Model", func() {
	Describe("Game Instance", func() {
		It("Should load a game", func() {
			//Given
			gameID := "game-id"
			var game models.Game

			//When
			err := db.
				Select("*").
				From("games").
				Where("id = $1", gameID).
				QueryStruct(&game)

			//Then
			Expect(err).NotTo(HaveOccurred())
			Expect(game.ID).To(Equal(gameID))
			Expect(game.Name).To(Equal("game-1"))

			obj, err := game.GetMetadata()
			Expect(err).NotTo(HaveOccurred())
			Expect(obj).To(BeEquivalentTo(map[string]interface{}{}))
		})

		It("Should create game", func() {
			//Given
			meta := dat.JSON([]byte(`{"qwe": 123}`))
			game := &models.Game{
				Name:     "Game Awesome Name",
				ID:       "game-id-2",
				Metadata: meta,
			}
			var game2 models.Game

			//When
			err1 := db.
				InsertInto("games").
				Columns("name", "id", "metadata").
				Record(game).
				Returning("created_at", "updated_at").
				QueryStruct(game)

			err3 := db.
				Select("*").
				From("games").
				Where("id = $1", game.ID).
				QueryStruct(&game2)
			obj, err2 := game2.GetMetadata()

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(err3).NotTo(HaveOccurred())
			Expect(game2.CreatedAt).NotTo(Equal(""))
			Expect(game2.ID).To(Equal(game.ID))
			Expect(game2.Name).To(Equal(game.Name))
			Expect(obj["qwe"]).To(BeEquivalentTo(123))
		})
	})

	Describe("List Games", func() {
		It("Should return the full list of games", func() {
			games, err := models.ListGames(db, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(games).To(HaveLen(7))
		})
	})

	Describe("Get game by id", func() {
		It("Should load game by id", func() {
			//Given
			gameID := "game-id"

			//When
			game, err1 := models.GetGameByID(db, gameID, nil)
			obj, err2 := game.GetMetadata()

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(game.ID).To(Equal(gameID))
			Expect(game.Name).To(Equal("game-1"))
			Expect(err2).NotTo(HaveOccurred())
			Expect(obj).To(BeEquivalentTo(map[string]interface{}{}))
		})

		It("Should return error if game not found", func() {
			//Given
			gameID := uuid.NewV4().String()
			expectedError := errors.NewModelNotFoundError("Game", map[string]interface{}{
				"ID": gameID,
			})

			//When
			game, err := models.GetGameByID(db, gameID, nil)

			//Then
			Expect(game.ID).To(Equal(""))
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(expectedError))
		})
	})

	Describe("Upsert game", func() {
		It("should insert game with new id", func() {
			//Given
			id := uuid.NewV4().String()
			game := models.Game{
				ID:       id,
				Name:     "Game Awesome Name",
				Metadata: dat.JSON([]byte(`{"qwe": 123}`)),
			}
			var c models.RealClock

			//When
			err1 := models.UpsertGame(db, &game, c.GetTime(), nil)
			gameFromDB, err2 := models.GetGameByID(db, id, nil)

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(gameFromDB.ID).To(Equal(id))
			Expect(gameFromDB.Name).To(Equal(game.Name))

			obj, err := gameFromDB.GetMetadata()
			Expect(err).NotTo(HaveOccurred())
			Expect(obj["qwe"]).To(BeEquivalentTo(123))
		})

		It("should update game with existing id", func() {
			//Given
			id := "upsert-game-id"
			name := "Game Awesome Name"
			game := models.Game{
				ID:   id,
				Name: name,
			}
			var c models.RealClock

			//When
			err1 := models.UpsertGame(db, &game, c.GetTime(), nil)
			gameFromDB, err2 := models.GetGameByID(db, id, nil)

			//Then
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(gameFromDB.ID).To(Equal(id))
			Expect(gameFromDB.Name).To(Equal(name))
		})

		It("should return error when inserting game with very big name", func() {
			//Given
			id := uuid.NewV4().String()

			var buffer bytes.Buffer
			for i := 0; i < 100; i++ {
				buffer.WriteString("abc")
			}

			game := models.Game{
				ID:       id,
				Name:     buffer.String(),
				Metadata: dat.JSON([]byte(`{"qwe": 123}`)),
			}
			var c models.RealClock

			//When
			err := models.UpsertGame(db, &game, c.GetTime(), nil)

			//Then
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("pq: value too long for type character varying(255)"))
		})
	})
})
