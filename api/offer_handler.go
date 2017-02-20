// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/topfreegames/offers/models"
)

//OfferRequestHandler handler
type OfferRequestHandler struct {
	App    *App
	Method string
}

//ServeHTTP method
func (h *OfferRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch h.Method {
	case "get_offers":
		h.getOffers(w, r)
	case "claim_offer":
		h.claimOffer(w, r)
	case "insert_offer":
		h.insertOffer(w, r)
	default:
		msg := "method not found"
		h.App.HandleError(w, http.StatusBadRequest, msg, errors.New(msg))
	}
}

func (h *OfferRequestHandler) getOffers(w http.ResponseWriter, r *http.Request) {
	//mr := metricsReporterFromCtx(r.Context())
	playerID := r.URL.Query().Get("player-id")
	if playerID == "" {
		err := fmt.Errorf("The player-id parameter cannot be empty")
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

func (h *OfferRequestHandler) claimOffer(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	offer := offerFromCtx(r.Context())
	currentTime := h.App.Clock.GetTime()

	err := models.ClaimOffer(h.App.DB, offer.ID, offer.GameID, currentTime, mr)

	if err != nil {
		h.App.HandleError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	Write(w, http.StatusOK, offer.ClaimedAt.Time.String())
}

func (h *OfferRequestHandler) insertOffer(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	offer := offerFromCtx(r.Context())
	currentTime := h.App.Clock.GetTime()

	err := models.UpsertOffer(h.App.DB, offer, currentTime, mr)

	if err != nil {
		h.App.HandleError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	Write(w, http.StatusOK, offer.ClaimedAt.Time.String())
}
