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
	"github.com/topfreegames/offers/models"
)

var _ = Describe("Games Model", func() {
	Describe("Game Instance", func() {
		It("Shoud load a game", func() {
			var game models.Game
			err := db.
				Select("*").
				From("games").
				Where("id = $1", "3393cd15-5b5a-4cfb-9725-ddbde660a727").
				QueryStruct(&game)

			Expect(err).NotTo(HaveOccurred())
			Expect(game.ID).To(Equal("3393cd15-5b5a-4cfb-9725-ddbde660a727"))
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

			game := models.Game{
				Name:     uuid.NewV4().String(),
				Metadata: meta,
			}
			err := db.
				InsertInto("games").
				Columns("name", "metadata").
				Record(&game).
				Returning("id", "created_at", "updated_at").
				QueryStruct(&game)

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
			game, err := models.GetGameByID(db, "3393cd15-5b5a-4cfb-9725-ddbde660a727")
			Expect(err).NotTo(HaveOccurred())
			Expect(game.ID).To(Equal("3393cd15-5b5a-4cfb-9725-ddbde660a727"))
			Expect(game.Name).To(Equal("game-1"))

			obj, err := game.GetMetadata()
			Expect(err).NotTo(HaveOccurred())
			Expect(obj).To(BeEquivalentTo(map[string]interface{}{}))
		})
	})
})
