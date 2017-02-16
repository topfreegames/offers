// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"net/http"

	"github.com/topfreegames/offers/models"
)

//GameHandler handler
type GameHandler struct {
	App *App
}

//ServeHTTP method
func (g *GameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	game := gameFromCtx(r.Context())

	err := mr.WithSegment(models.SegmentModel, func() error {
		return models.UpsertGame(g.App.DB, game, mr)
	})

	if err != nil {
		Write(w, http.StatusBadRequest, "Creating game failed.")
		return
	}

	Write(w, http.StatusOK, game.ID)
}
