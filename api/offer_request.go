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

	"github.com/sirupsen/logrus"
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
		h.viewOffer(w, r)
	case "offer-info":
		h.offerInfo(w, r)
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
		h.App.HandleError(w, http.StatusBadRequest, err.Error(), err)
		return
	} else if gameID == "" {
		err := fmt.Errorf("The game-id parameter cannot be empty")
		logger.WithError(err).Error("Failed to retrieve offer for player.")
		h.App.HandleError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	currentTime := h.App.Clock.GetTime()
	filterAttrsList := r.URL.Query()
	filterAttrs := make(map[string]string)
	delete(filterAttrsList, "player-id")
	delete(filterAttrsList, "game-id")
	for k, v := range filterAttrsList {
		if len(v) == 0 || len(v) > 1 {
			err := fmt.Errorf("Filter attribute passed with invalid number of arguments. Key: %s", k)
			logger.WithError(err).Error("Failed to retrieve offer for player.")
			h.App.HandleError(w, http.StatusBadRequest, "A filter parameter is invalid.", err)
			return
		}
		filterAttrs[k] = v[0]
	}

	maxAge := h.App.MaxAge
	allowInefficientQueries := false
	var game *models.Game
	var err error
	err = mr.WithSegment(models.SegmentModel, func() error {
		game, err = models.GetGameByID(r.Context(), h.App.DB, gameID, mr)
		return err
	})
	if err != nil {
		logger.WithError(err).Error("Failed to retrieve game.")
		if modelNotFound, ok := err.(*e.ModelNotFoundError); ok {
			h.App.HandleError(w, http.StatusNotFound, modelNotFound.Error(), modelNotFound)
			return
		}
		h.App.HandleError(w, http.StatusInternalServerError, "Failed to retrieve game", err)
		return
	}
	metadata, err := game.GetMetadata()
	if err != nil {
		logger.WithError(err).Error("Failed to get game metadata.")
		h.App.HandleError(w, http.StatusInternalServerError, "Failed to get game metadata", err)
		return
	}
	if val, ok := metadata["allowInefficientQueries"]; ok {
		if allowFromMeta, ok := val.(bool); ok {
			allowInefficientQueries = allowFromMeta
		}
	}
	if val, ok := metadata["cacheMaxAge"]; ok {
		if maxAgeFromMeta, ok := val.(float64); ok {
			maxAge = int64(maxAgeFromMeta)
		}
	}

	var offers map[string][]*models.OfferToReturn
	err = mr.WithSegment(models.SegmentModel, func() error {
		offers, err = models.GetAvailableOffers(r.Context(), h.App.DB, h.App.Cache, gameID, playerID, currentTime, h.App.OffersCacheMaxAge, filterAttrs, allowInefficientQueries, mr)
		return err
	})

	if err != nil {
		logger.WithError(err).Error("Failed to retrieve offers for player.")
		h.App.HandleError(w, http.StatusInternalServerError, "Failed to retrieve offer for player", err)
		return
	}

	bytes, err := json.Marshal(offers)
	if err != nil {
		logger.WithError(err).Error("Failed to parse structs to JSON.")
		h.App.HandleError(w, http.StatusInternalServerError, "Failed to parse structs to JSON", err)
		return
	}

	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", maxAge))

	logger.Debug("Retrieved player offers successfully")
	WriteBytes(w, http.StatusOK, bytes)
}

func (h *OfferRequestHandler) claimOffer(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	payload := claimOfferPayloadFromCtx(r.Context())
	currentTime := h.App.Clock.GetTime()
	logger := h.App.Logger.WithFields(logrus.Fields{
		"source":    "offerHandler",
		"operation": "claimOffer",
		"payload":   payload,
	})

	if payload.OfferInstanceID == "" && payload.ProductID == "" {
		validationError := e.NewValidationFailedError(errors.New("Id and ProductID cannot be both null"))
		h.App.HandleError(w, http.StatusUnprocessableEntity, validationError.Error(), validationError)
		return
	}

	contents, alreadyClaimed, nextAt, err := models.ClaimOffer(
		r.Context(),
		h.App.DB,
		payload.GameID,
		payload.OfferInstanceID,
		payload.PlayerID,
		payload.ProductID,
		payload.TransactionID,
		payload.Timestamp,
		currentTime,
		mr,
	)

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

func (h *OfferRequestHandler) viewOffer(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	payload := offerImpressionPayloadFromCtx(r.Context())
	offerInstanceID := paramKeyFromContext(r.Context())
	currentTime := h.App.Clock.GetTime()
	logger := h.App.Logger.WithFields(logrus.Fields{
		"source":    "offerHandler",
		"operation": "viewOffer",
		"payload":   payload,
	})

	alreadyViewed, nextAt, err := models.ViewOffer(
		r.Context(),
		h.App.DB,
		payload.GameID,
		offerInstanceID,
		payload.PlayerID,
		payload.ImpressionID,
		currentTime,
		mr,
	)
	if err != nil {
		logger.WithError(err).Error("Failed to update offer impressions.")
		if modelNotFound, ok := err.(*e.ModelNotFoundError); ok {
			h.App.HandleError(w, http.StatusNotFound, modelNotFound.Error(), modelNotFound)
			return
		}

		h.App.HandleError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	logger.Debug("Updated offer impressions successfully")
	res := map[string]interface{}{}
	if nextAt != 0 {
		res["nextAt"] = nextAt
	}
	bytesRes, _ := json.Marshal(res)
	if alreadyViewed {
		WriteBytes(w, http.StatusConflict, bytesRes)
		return
	}
	WriteBytes(w, http.StatusOK, bytesRes)
}

func (h *OfferRequestHandler) offerInfo(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	playerID := r.URL.Query().Get("player-id")
	gameID := r.URL.Query().Get("game-id")
	offerInstanceID := r.URL.Query().Get("offer-id")
	logger := h.App.Logger.WithFields(logrus.Fields{
		"source":          "offerHandler",
		"operation":       "offerInfo",
		"gameID":          gameID,
		"playerID":        playerID,
		"offerInstanceID": offerInstanceID,
	})

	if playerID == "" {
		err := fmt.Errorf("The player-id parameter cannot be empty")
		logger.WithError(err).Error("Failed to retrieve offer for player.")
		h.App.HandleError(w, http.StatusBadRequest, err.Error(), err)
		return
	} else if gameID == "" {
		err := fmt.Errorf("The game-id parameter cannot be empty")
		logger.WithError(err).Error("Failed to retrieve offer for player.")
		h.App.HandleError(w, http.StatusBadRequest, err.Error(), err)
		return
	} else if offerInstanceID == "" {
		err := fmt.Errorf("The offer-id parameter cannot be empty")
		logger.WithError(err).Error("Failed to retrieve offer for player.")
		h.App.HandleError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	var err error
	var offer *models.OfferToReturn
	err = mr.WithSegment(models.SegmentModel, func() error {
		offer, err = models.GetOfferInfo(r.Context(), h.App.DB, gameID, offerInstanceID, h.App.OffersCacheMaxAge, mr)
		return err
	})

	if err != nil {
		if modelNotFound, ok := err.(*e.ModelNotFoundError); ok {
			logger.WithError(err).Error("Failed to find offer for player.")
			h.App.HandleError(w, http.StatusNotFound, modelNotFound.Error(), modelNotFound)
			return
		}
		logger.WithError(err).Error("Failed to retrieve offer for player.")
		h.App.HandleError(w, http.StatusInternalServerError, "Failed to retrieve offer for player", err)
		return
	}

	bytes, err := json.Marshal(offer)
	if err != nil {
		logger.WithError(err).Error("Failed to parse structs to JSON.")
		h.App.HandleError(w, http.StatusInternalServerError, "Failed to parse structs to JSON", err)
		return
	}

	maxAge := h.App.MaxAge
	var game *models.Game
	err = mr.WithSegment(models.SegmentModel, func() error {
		game, err = models.GetGameByID(r.Context(), h.App.DB, gameID, mr)
		return err
	})
	metadata, err := game.GetMetadata()
	if err == nil {
		if val, ok := metadata["cacheMaxAge"]; ok {
			if maxAgeFromMeta, ok := val.(float64); ok {
				maxAge = int64(maxAgeFromMeta)
			}
		}
	} else {
		logger.WithError(err).Warn("Failed to get game metadata.")
	}
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", maxAge))

	logger.Info("Retrieved player offer successfully")
	WriteBytes(w, http.StatusOK, bytes)
}
