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

	e "github.com/topfreegames/offers/errors"
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
	case "update_offer_last_seen_at":
		h.updateOfferLastSeenAt(w, r)
	}
}

func (h *OfferRequestHandler) getOffers(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	playerID := r.URL.Query().Get("player-id")
	gameID := r.URL.Query().Get("game-id")
	if playerID == "" {
		err := fmt.Errorf("The player-id parameter cannot be empty")
		h.App.HandleError(w, http.StatusBadRequest, "The player-id parameter cannot be empty.", err)
		return
	} else if gameID == "" {
		err := fmt.Errorf("The game-id parameter cannot be empty")
		h.App.HandleError(w, http.StatusBadRequest, "The game-id parameter cannot be empty.", err)
		return
	}
	currentTime := h.App.Clock.GetTime()

	ots, err := models.GetAvailableOffers(h.App.DB, playerID, gameID, currentTime, mr)

	if err != nil {
		h.App.HandleError(w, http.StatusInternalServerError, "Failed to retrieve offer for player", err)
		return
	}

	bytes, err := json.Marshal(ots)

	if err != nil {
		h.App.HandleError(w, http.StatusInternalServerError, "Failed to parse structs to JSON", err)
		return
	}

	WriteBytes(w, http.StatusOK, bytes)
}

func (h *OfferRequestHandler) claimOffer(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	offer := offerToUpdateFromCtx(r.Context())
	currentTime := h.App.Clock.GetTime()

	contents, alreadyClaimed, err := models.ClaimOffer(h.App.DB, offer.ID, offer.PlayerID, offer.GameID, currentTime, mr)

	if err != nil {
		if modelNotFound, ok := err.(*e.ModelNotFoundError); ok {
			h.App.HandleError(w, http.StatusNotFound, modelNotFound.Error(), modelNotFound)
			return
		}

		h.App.HandleError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	if alreadyClaimed {
		WriteBytes(w, http.StatusConflict, contents)
		return
	}

	WriteBytes(w, http.StatusOK, contents)
}

//UpdateOfferLastSeenAt updates the offer last seen at
func (h *OfferRequestHandler) updateOfferLastSeenAt(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	offer := offerToUpdateFromCtx(r.Context())
	currentTime := h.App.Clock.GetTime()

	err := models.UpdateOfferLastSeenAt(h.App.DB, offer.ID, offer.PlayerID, offer.GameID, currentTime, mr)

	if err != nil {
		if modelNotFound, ok := err.(*e.ModelNotFoundError); ok {
			h.App.HandleError(w, http.StatusNotFound, modelNotFound.Error(), modelNotFound)
			return
		}

		h.App.HandleError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	Write(w, http.StatusOK, "")
}
