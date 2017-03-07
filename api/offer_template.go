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
	"github.com/topfreegames/offers/errors"
	"github.com/topfreegames/offers/models"
)

//OfferTemplateHandler handler
type OfferTemplateHandler struct {
	App    *App
	Method string
}

//ServeHTTP method
func (g *OfferTemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch g.Method {
	case "insert":
		g.insert(w, r, true)
		return
	case "enable":
		g.setEnabledOfferTemplate(w, r, true)
	case "disable":
		g.setEnabledOfferTemplate(w, r, false)
		return
	case "list":
		g.list(w, r)
		return
	case "update":
		g.insert(w, r, false)
		return
	}
}

func (g *OfferTemplateHandler) insert(w http.ResponseWriter, r *http.Request, onlyInsert bool) {
	mr := metricsReporterFromCtx(r.Context())
	ot := offerTemplateFromCtx(r.Context())
	offerTemplateKey := paramKeyFromContext(r.Context())
	ot.Key = offerTemplateKey
	userEmail := userEmailFromContext(r.Context())

	logger := g.App.Logger.WithFields(logrus.Fields{
		"source":        "offerTemplateHandler",
		"operation":     "insertOfferTemplate",
		"userEmail":     userEmail,
		"offerTemplate": ot,
	})

	var err error

	err = mr.WithSegment(models.SegmentModel, func() error {
		ot, err = models.InsertOfferTemplate(g.App.DB, ot, onlyInsert, mr)
		return err
	})

	if err != nil {
		logger.WithError(err).Error("Insert offer template failed.")
		if foreignKeyError, ok := err.(*errors.InvalidModelError); ok {
			g.App.HandleError(w, http.StatusUnprocessableEntity, foreignKeyError.Error(), foreignKeyError)
			return
		}

		if conflictedKeyError, ok := err.(*errors.ConflictedModelError); ok {
			g.App.HandleError(w, http.StatusConflict, conflictedKeyError.Error(), conflictedKeyError)
			return
		}

		g.App.HandleError(w, http.StatusInternalServerError, "Insert offer template failed", err)
		return
	}

	bytesRes, err := json.Marshal(ot)
	if err != nil {
		logger.WithError(err).Error("Failed to build offer template response.")
		g.App.HandleError(w, http.StatusInternalServerError, "Failed to build offer template response", err)
		return
	}

	logger.Info("Inserted offer template successfuly.")
	WriteBytes(w, http.StatusCreated, bytesRes)
}

func (g *OfferTemplateHandler) update(w http.ResponseWriter, r *http.Request) {
}

func (g *OfferTemplateHandler) setEnabledOfferTemplate(w http.ResponseWriter, r *http.Request, enable bool) {
	mr := metricsReporterFromCtx(r.Context())
	offerTemplateID := paramKeyFromContext(r.Context())
	userEmail := userEmailFromContext(r.Context())

	logger := g.App.Logger.WithFields(logrus.Fields{
		"source":          "offerTemplateHandler",
		"operation":       "setEnabledOfferTemplate",
		"userEmail":       userEmail,
		"offerTemplateID": offerTemplateID,
	})

	var err error
	err = mr.WithSegment(models.SegmentModel, func() error {
		return models.SetEnabledOfferTemplate(g.App.DB, offerTemplateID, enable, mr)
	})

	if err != nil {
		logger.WithError(err).Error("Update offer template failed.")
		if modelNotFound, ok := err.(*errors.ModelNotFoundError); ok {
			g.App.HandleError(w, http.StatusNotFound, "Offer template not found for this ID", modelNotFound)
			return
		}
		g.App.HandleError(w, http.StatusInternalServerError, "Update offer template failed", err)
		return
	}

	logger.Info("Updated offer template successfuly.")
	bytesRes, _ := json.Marshal(map[string]interface{}{"id": offerTemplateID})
	WriteBytes(w, http.StatusOK, bytesRes)
}

func (g *OfferTemplateHandler) list(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	gameID := r.URL.Query().Get("game-id")
	userEmail := userEmailFromContext(r.Context())

	logger := g.App.Logger.WithFields(logrus.Fields{
		"source":    "offerTemplateHandler",
		"operation": "list",
		"userEmail": userEmail,
		"gameID":    gameID,
	})

	if gameID == "" {
		err := fmt.Errorf("The game-id parameter cannot be empty")
		logger.WithError(err).Error("List game offer templates failed.")
		g.App.HandleError(w, http.StatusBadRequest, "The game-id parameter cannot be empty.", err)
		return
	}

	var err error
	var offerTemplates []*models.OfferTemplate
	err = mr.WithSegment(models.SegmentModel, func() error {
		offerTemplates, err = models.ListOfferTemplates(g.App.DB, gameID, mr)
		return err
	})

	if err != nil {
		logger.WithError(err).Error("List game offer templates failed.")
		g.App.HandleError(w, http.StatusInternalServerError, "List game offer templates failed.", err)
		return
	}

	logger.Info("Listed game offer templates successfully.")
	if len(offerTemplates) == 0 {
		Write(w, http.StatusOK, "[]")
		return
	}
	bytes, _ := json.Marshal(offerTemplates)
	WriteBytes(w, http.StatusOK, bytes)
}
