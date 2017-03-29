// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/topfreegames/offers/models"
)

//GameHandler handler
type GameHandler struct {
	App    *App
	Method string
}

//ServeHTTP method
func (g *GameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch g.Method {
	case "list":
		g.list(w, r)
		return
	case "upsert":
		g.upsert(w, r)
		return
	}
}

func (g *GameHandler) list(w http.ResponseWriter, r *http.Request) {
	userEmail := userEmailFromContext(r.Context())
	logger := g.App.Logger.WithFields(logrus.Fields{
		"source":    "gameHandler",
		"operation": "list",
		"userEmail": userEmail,
	})
	mr := metricsReporterFromCtx(r.Context())

	var err error
	var games []*models.Game
	err = mr.WithSegment(models.SegmentModel, func() error {
		games, err = models.ListGames(g.App.DB, mr)
		return err
	})

	if err != nil {
		logger.WithError(err).Error("List games failed.")
		g.App.HandleError(w, http.StatusInternalServerError, "List games failed.", err)
		return
	}

	logger.Info("Listed games successfully.")
	if len(games) == 0 {
		Write(w, http.StatusOK, "[]")
		return
	}
	bytes, _ := json.Marshal(games)
	WriteBytes(w, http.StatusOK, bytes)
}

func (g *GameHandler) upsert(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	game := gameFromCtx(r.Context())
	userEmail := userEmailFromContext(r.Context())

	logger := g.App.Logger.WithFields(logrus.Fields{
		"source":    "gameHandler",
		"operation": "upsert",
		"userEmail": userEmail,
		"game":      game,
	})

	err := mr.WithSegment(models.SegmentModel, func() error {
		currentTime := g.App.Clock.GetTime()
		return models.UpsertGame(g.App.DB, game, currentTime, mr)
	})

	if err != nil {
		logger.WithError(err).Error("Upserting game failed.")
		g.App.HandleError(w, http.StatusInternalServerError, "Upserting game failed", err)
		return
	}
	logger.Info("Upserted game successfully.")
	bytesRes, _ := json.Marshal(map[string]interface{}{"gameId": game.ID})
	WriteBytes(w, http.StatusOK, bytesRes)
}
