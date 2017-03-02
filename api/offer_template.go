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
		g.insertOfferTemplate(w, r)
		return
	case "set-enabled":
		g.setEnabledOfferTemplate(w, r)
		return
	}
}

func (g *OfferTemplateHandler) insertOfferTemplate(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	ot := offerTemplateFromCtx(r.Context())

	fmt.Printf("OFFER TEMPLATE %#v", ot)

	var err error
	err = mr.WithSegment(models.SegmentModel, func() error {
		ot, err = models.InsertOfferTemplate(g.App.DB, ot, mr)
		return err
	})

	if err != nil {
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
		g.App.HandleError(w, http.StatusInternalServerError, "Failed to build offer template response", err)
		return
	}

	WriteBytes(w, http.StatusOK, bytesRes)
}

func (g *OfferTemplateHandler) setEnabledOfferTemplate(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	ot := offerTemplateToUpdateFromCtx(r.Context())

	var err error
	err = mr.WithSegment(models.SegmentModel, func() error {
		return models.SetEnabledOfferTemplate(g.App.DB, ot.ID, ot.Enabled, mr)
	})

	if err != nil {
		if modelNotFound, ok := err.(*errors.ModelNotFoundError); ok {
			g.App.HandleError(w, http.StatusNotFound, "Offer template not found for this ID", modelNotFound)
			return
		}
		g.App.HandleError(w, http.StatusInternalServerError, "Update offer template failed", err)
		return
	}

	Write(w, http.StatusOK, ot.ID)
}
