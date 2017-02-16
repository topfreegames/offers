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

	"github.com/topfreegames/offers/models"
)

//GameHandler handler
type GameHandler struct {
	App *App
}

//ServeHTTP method
func (g *GameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var game models.Game
	err := decoder.Decode(&game)

	if err != nil {
		http.Error(w, "empty parameter", 400)
	} else if err = models.UpsertGame(g.App.DB, &game); err != nil {
		http.Error(w, "upsert error", 400)
	}
}
