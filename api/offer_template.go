// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"net/http"

	"github.com/topfreegames/offers/models"
)

//OfferTemplateHandler handler
type OfferTemplateHandler struct {
	App *App
}

//ServeHTTP method
func (g *OfferTemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mr := metricsReporterFromCtx(r.Context())
	ot := offerTemplateFromCtx(r.Context())

	err := mr.WithSegment(models.SegmentModel, func() error {
		return models.InsertOfferTemplate(g.App.DB, ot, mr)
	})

	if err != nil {
		g.App.HandleError(w, http.StatusInternalServerError, "Insert offer template failed", err)
		return
	}

	Write(w, http.StatusOK, ot.ID)
}
