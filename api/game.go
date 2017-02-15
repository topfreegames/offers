// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"fmt"
	"net/http"
)

//GameHandler handler
type GameHandler struct {
	App *App
}

//ServeHTTP method
func (g *GameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Printf("%v", r.URL)
	}
}
