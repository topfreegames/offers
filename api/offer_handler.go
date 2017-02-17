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
)

//OfferRequestHandler handler
type OfferRequestHandler struct {
	App *App
}

//ServeHTTP method
func (h *OfferRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//mr := metricsReporterFromCtx(r.Context())
	payload := offerRequestPayloadFromCtx(r.Context())
	playerID := payload.PlayerID
	currentTime := h.App.Clock.GetTime()

	l := loggerFromContext(r.Context()).WithFields(logrus.Fields{
		"playerID":    playerID,
		"currentTime": currentTime,
	})

	l.Debug("Retrieving player info...")

	//info := models.GetPlayerOffers(playerID)
	//availableOffers := models.GetAvailableOffers()

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
