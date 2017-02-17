// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
)

//OfferRequestHandler handler
type OfferRequestHandler struct {
	App *App
}

//ServeHTTP method
func (h *OfferRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//mr := metricsReporterFromCtx(r.Context())
	playerID := r.URL.Query().Get("player-id")
	if playerID == "" {
		err := fmt.Errorf("The player-id parameter cannot be empty.")
		h.App.HandleError(w, http.StatusBadRequest, "The player-id parameter cannot be empty.", err)
		return
	}
	currentTime := h.App.Clock.GetTime()

	l := loggerFromContext(r.Context()).WithFields(logrus.Fields{
		"playerID":    playerID,
		"currentTime": currentTime,
	})

	l.Debug("Retrieving player info...")

	//availableOffers, err := models.GetEnabledOfferTemplates(h.App.DB, mr)
	//if err != nil {
	//h.App.HandleError(w, http.StatusInternalServerError, "Failed to retrieve enabled offers", err)
	//return
	//}

	//info := models.GetPlayerSeenOffers(h.App.DB, playerID, availableOffers, mr)

	//fitOffers := []*OfferPayload
	//for _, offer := range availableOffers {
	//if offer.Fits(info) {
	//fitOffers = append(fitOffers, OfferPayloadFrom(offer))
	//}
	//}

	res, err := json.Marshal(map[string]interface{}{
	//"offers": fitOffers,
	})

	if err != nil {
		h.App.HandleError(w, http.StatusInternalServerError, "Failed to retrieve offer for player", err)
		return
	}

	WriteBytes(w, http.StatusOK, res)
	l.Debug("Offer request done.")
}
