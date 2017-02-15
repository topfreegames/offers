// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import "net/http"

//HealthcheckHandler handler
type HealthcheckHandler struct {
	App *App
}

//ServeHTTP method
func (h *HealthcheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := loggerFromContext(r.Context())

	l.Debug("Performing healthcheck...")

	_, err := h.App.DB.Exec("select 1")
	if err != nil {
		Write(w, http.StatusInternalServerError, "Database is offline")
		l.WithError(err).Error("Database is offline")
		return
	}

	Write(w, http.StatusOK, "WORKING")
	l.Debug("Healthcheck done.")
}
