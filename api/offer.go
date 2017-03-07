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
	case "get-offers":
		h.getOffers(w, r)
	case "claim":
		h.claimOffer(w, r)
	case "impressions":
		h.updateOfferLastSeenAt(w, r)
	}
}

func (h *OfferRequestHandler) getOffers(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	playerID := r.URL.Query().Get("player-id")
	gameID := r.URL.Query().Get("game-id")
	logger := h.App.Logger.WithFields(logrus.Fields{
		"source":    "offerHandler",
		"operation": "getOffers",
		"gameID":    gameID,
		"playerID":  playerID,
	})

	if playerID == "" {
		err := fmt.Errorf("The player-id parameter cannot be empty")
		logger.WithError(err).Error("Failed to retrieve offer for player.")
		h.App.HandleError(w, http.StatusBadRequest, "The player-id parameter cannot be empty.", err)
		return
	} else if gameID == "" {
		err := fmt.Errorf("The game-id parameter cannot be empty")
		logger.WithError(err).Error("Failed to retrieve offer for player.")
		h.App.HandleError(w, http.StatusBadRequest, "The game-id parameter cannot be empty.", err)
		return
	}
	currentTime := h.App.Clock.GetTime()

	ots, err := models.GetAvailableOffers(h.App.DB, playerID, gameID, currentTime, mr)

	if err != nil {
		logger.WithError(err).Error("Failed to retrieve offer for player.")
		h.App.HandleError(w, http.StatusInternalServerError, "Failed to retrieve offer for player", err)
		return
	}

	bytes, err := json.Marshal(ots)

	if err != nil {
		logger.WithError(err).Error("Failed to parse structs to JSON.")
		h.App.HandleError(w, http.StatusInternalServerError, "Failed to parse structs to JSON", err)
		return
	}

	logger.Info("Retrieved player offers successfully")
	WriteBytes(w, http.StatusOK, bytes)
}

func (h *OfferRequestHandler) claimOffer(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	offer := offerToUpdateFromCtx(r.Context())
	offerID := paramKeyFromContext(r.Context())
	currentTime := h.App.Clock.GetTime()
	logger := h.App.Logger.WithFields(logrus.Fields{
		"source":    "offerHandler",
		"operation": "claimOffer",
		"offer":     offer,
		"offerID":   offerID,
	})

	contents, alreadyClaimed, nextAt, err := models.ClaimOffer(h.App.DB, offerID, offer.PlayerID, offer.GameID, currentTime, mr)
	if err != nil {
		logger.WithError(err).Error("Failed to claim offer.")
		if modelNotFound, ok := err.(*e.ModelNotFoundError); ok {
			h.App.HandleError(w, http.StatusNotFound, modelNotFound.Error(), modelNotFound)
			return
		}

		h.App.HandleError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	logger.WithField("alreadyClaimed", alreadyClaimed).Info("Claimed offer successfully")
	res := map[string]interface{}{
		"contents": contents,
	}
	if nextAt != 0 {
		res["nextAt"] = nextAt
	}
	bytesRes, _ := json.Marshal(res)
	if alreadyClaimed {
		WriteBytes(w, http.StatusConflict, bytesRes)
		return
	}

	WriteBytes(w, http.StatusOK, bytesRes)
}

//UpdateOfferLastSeenAt updates the offer last seen at
func (h *OfferRequestHandler) updateOfferLastSeenAt(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	offer := offerToUpdateFromCtx(r.Context())
	offerID := paramKeyFromContext(r.Context())
	currentTime := h.App.Clock.GetTime()
	logger := h.App.Logger.WithFields(logrus.Fields{
		"source":    "offerHandler",
		"operation": "updateOfferLastSeenAt",
		"offer":     offer,
		"offerID":   offerID,
	})

	nextAt, err := models.UpdateOfferLastSeenAt(h.App.DB, offerID, offer.PlayerID, offer.GameID, currentTime, mr)
	if err != nil {
		logger.WithError(err).Error("Failed to updated offer impressions.")
		if modelNotFound, ok := err.(*e.ModelNotFoundError); ok {
			h.App.HandleError(w, http.StatusNotFound, modelNotFound.Error(), modelNotFound)
			return
		}

		h.App.HandleError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	logger.Info("Upated offer impressions successfully")
	if nextAt == 0 {
		Write(w, http.StatusOK, "{}")
		return
	}
	bytesRes, _ := json.Marshal(map[string]interface{}{"nextAt": nextAt})
	WriteBytes(w, http.StatusOK, bytesRes)
}
