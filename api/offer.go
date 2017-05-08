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
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/topfreegames/offers/errors"
	"github.com/topfreegames/offers/models"
)

//OfferHandler handler
type OfferHandler struct {
	App    *App
	Method string
}

//ServeHTTP method
func (g *OfferHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch g.Method {
	case "insert":
		g.insertOffer(w, r)
		return
	case "update":
		g.updateOffer(w, r)
		return
	case "enable":
		g.setEnabledOffer(w, r, true)
		return
	case "disable":
		g.setEnabledOffer(w, r, false)
		return
	case "list":
		g.list(w, r)
		return
	}
}

func (g *OfferHandler) insertOffer(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	offer := offerFromCtx(r.Context())
	userEmail := userEmailFromContext(r.Context())

	logger := g.App.Logger.WithFields(logrus.Fields{
		"source":    "offerHandler",
		"operation": "insertOffer",
		"userEmail": userEmail,
		"offer":     offer,
	})

	var err error
	err = mr.WithSegment(models.SegmentModel, func() error {
		offer, err = models.InsertOffer(g.App.DB, offer, mr)
		return err
	})

	if err != nil {
		logger.WithError(err).Error("Insert offer failed.")
		if foreignKeyError, ok := err.(*errors.InvalidModelError); ok {
			g.App.HandleError(w, http.StatusUnprocessableEntity, foreignKeyError.Error(), foreignKeyError)
			return
		}

		g.App.HandleError(w, http.StatusInternalServerError, "Insert offer failed", err)
		return
	}

	bytesRes, err := json.Marshal(offer)
	if err != nil {
		logger.WithError(err).Error("Failed to build offer response.")
		g.App.HandleError(w, http.StatusInternalServerError, "Failed to build offer response", err)
		return
	}

	logger.Info("Inserted offer successfuly.")
	WriteBytes(w, http.StatusCreated, bytesRes)
}

func (g *OfferHandler) setEnabledOffer(w http.ResponseWriter, r *http.Request, enable bool) {
	mr := metricsReporterFromCtx(r.Context())
	offerID := paramKeyFromContext(r.Context())
	userEmail := userEmailFromContext(r.Context())
	gameID := r.URL.Query().Get("game-id")

	logger := g.App.Logger.WithFields(logrus.Fields{
		"source":    "offerHandler",
		"operation": "setEnabledOffer",
		"userEmail": userEmail,
		"offerID":   offerID,
		"gameID":    gameID,
	})

	var err error
	err = mr.WithSegment(models.SegmentModel, func() error {
		return models.SetEnabledOffer(g.App.DB, gameID, offerID, enable, mr)
	})

	if err != nil {
		logger.WithError(err).Error("Update offer failed.")
		if modelNotFound, ok := err.(*errors.ModelNotFoundError); ok {
			g.App.HandleError(w, http.StatusNotFound, "Offer not found for this ID", modelNotFound)
			return
		}
		g.App.HandleError(w, http.StatusInternalServerError, "Update offer failed", err)
		return
	}

	logger.Info("Updated offer successfuly.")
	bytesRes, _ := json.Marshal(map[string]interface{}{"id": offerID})
	WriteBytes(w, http.StatusOK, bytesRes)
}

func (g *OfferHandler) updateOffer(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	offer := offerFromCtx(r.Context())
	userEmail := userEmailFromContext(r.Context())
	offerID := paramKeyFromContext(r.Context())
	offer.ID = offerID
	logger := g.App.Logger.WithFields(logrus.Fields{
		"source":    "offerHandler",
		"operation": "updateOffer",
		"userEmail": userEmail,
		"offer":     offer,
	})

	var err error
	err = mr.WithSegment(models.SegmentModel, func() error {
		offer, err = models.UpdateOffer(g.App.DB, offer, mr)
		return err
	})
	if err != nil {
		logger.WithError(err).Error("Update offer failed.")
		if notFoundError, ok := err.(*errors.ModelNotFoundError); ok {
			g.App.HandleError(w, http.StatusNotFound, notFoundError.Error(), notFoundError)
			return
		}

		g.App.HandleError(w, http.StatusInternalServerError, "Update offer failed", err)
		return
	}

	bytesRes, err := json.Marshal(offer)
	if err != nil {
		logger.WithError(err).Error("Failed to build offer response.")
		g.App.HandleError(w, http.StatusInternalServerError, "Failed to build offer response", err)
		return
	}

	logger.Info("Updated offer successfuly.")
	WriteBytes(w, http.StatusOK, bytesRes)
}

func (g *OfferHandler) list(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())

	vars := r.URL.Query()
	gameID := vars.Get("game-id")
	limitStr := vars.Get("limit")
	offsetStr := vars.Get("offset")
	var err error
	userEmail := userEmailFromContext(r.Context())

	logger := g.App.Logger.WithFields(logrus.Fields{
		"source":    "offerHandler",
		"operation": "list",
		"userEmail": userEmail,
		"gameID":    gameID,
	})

	if gameID == "" {
		err := fmt.Errorf("The game-id parameter cannot be empty")
		logger.WithError(err).Error("List game offers failed.")
		g.App.HandleError(w, http.StatusBadRequest, "The game-id parameter cannot be empty.", err)
		return
	}

	var limit uint64
	if limitStr == "" {
		limit = g.App.Pagination.Limit
	} else {
		var err error
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			logger.WithError(err).Error("List game offers failed.")
			g.App.HandleError(w, http.StatusBadRequest, "The limit parameter must be an uint.", err)
			return
		}
	}

	var offset uint64
	if offsetStr == "" {
		offset = g.App.Pagination.Offset
	} else {
		offset, err = strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			logger.WithError(err).Error("List game offers failed.")
			g.App.HandleError(w, http.StatusBadRequest, "The offset parameter must be an uint.", err)
			return
		}
	}

	var offers []*models.Offer
	var pages int
	err = mr.WithSegment(models.SegmentModel, func() error {
		offers, pages, err = models.ListOffers(g.App.DB, gameID, limit, offset, mr)
		return err
	})

	if err != nil {
		logger.WithError(err).Error("List game offers failed.")
		g.App.HandleError(w, http.StatusInternalServerError, "List game offers failed.", err)
		return
	}

	logger.Info("Listed game offers successfully.")
	responseObj := map[string]interface{}{
		"offers": offers,
		"pages":  pages,
	}

	bts, _ := json.Marshal(responseObj)
	WriteBytes(w, http.StatusOK, bts)
}
