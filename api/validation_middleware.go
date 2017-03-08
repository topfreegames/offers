// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"context"
	"encoding/json"
	"net/http"

	"gopkg.in/mgutz/dat.v2/dat"

	"github.com/asaskevich/govalidator"
	"github.com/topfreegames/offers/errors"
	"github.com/topfreegames/offers/models"
)

//ValidationMiddleware adds the version to the request
type ValidationMiddleware struct {
	GetPayload func() interface{}
	next       http.Handler
}

type contextKey string

const payloadString = contextKey("payload")

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware(f func() interface{}) *ValidationMiddleware {
	m := &ValidationMiddleware{GetPayload: f}
	m.configureCustomValidators()
	return m
}

func newContextWithPayload(ctx context.Context, payload interface{}, r *http.Request) context.Context {
	c := context.WithValue(ctx, payloadString, payload)
	return c
}

func gameFromCtx(ctx context.Context) *models.Game {
	game := ctx.Value(payloadString)
	if game == nil {
		return nil
	}

	gameWithID := game.(*models.Game)
	gameWithID.ID = paramKeyFromContext(ctx)

	return gameWithID
}

func offerFromCtx(ctx context.Context) *models.Offer {
	offer := ctx.Value(payloadString)
	if offer == nil {
		return nil
	}
	return offer.(*models.Offer)
}

func offerImpressionPayloadFromCtx(ctx context.Context) *models.OfferImpressionPayload {
	payload := ctx.Value(payloadString)
	if payload == nil {
		return nil
	}
	return payload.(*models.OfferImpressionPayload)
}

func claimOfferPayloadFromCtx(ctx context.Context) *models.ClaimOfferPayload {
	payload := ctx.Value(payloadString)
	if payload == nil {
		return nil
	}
	return payload.(*models.ClaimOfferPayload)
}

func (m *ValidationMiddleware) configureCustomValidators() {
	govalidator.CustomTypeTagMap.Set(
		"RequiredJSONObject",
		govalidator.CustomTypeValidator(
			func(i interface{}, context interface{}) bool {
				switch v := i.(type) {
				case dat.JSON:
					var val map[string]interface{}
					err := v.Unmarshal(&val)
					return err == nil && len(val) > 0
				}
				return false
			},
		),
	)
	govalidator.CustomTypeTagMap.Set(
		"JSONObject",
		govalidator.CustomTypeValidator(
			func(i interface{}, context interface{}) bool {
				switch v := i.(type) {
				case dat.JSON:
					var val map[string]interface{}
					err := v.Unmarshal(&val)
					if err == nil {
						return true
					}
					m, err := v.MarshalJSON()
					return err == nil && string(m) == "null"
				}
				return false
			},
		),
	)
}

//ServeHTTP method
func (m *ValidationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	payload := m.GetPayload()
	l := loggerFromContext(r.Context())

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(payload)

	if err != nil {
		l.WithError(err).Error("Payload could not be decoded.")
		vErr := errors.NewValidationFailedError(err)
		WriteBytes(w, http.StatusBadRequest, vErr.Serialize())
		return
	}

	_, err = govalidator.ValidateStruct(payload)

	if err != nil {
		l.WithError(err).Error("Payload is invalid.")
		vErr := errors.NewValidationFailedError(err)
		WriteBytes(w, http.StatusUnprocessableEntity, vErr.Serialize())
		return
	}

	c := newContextWithPayload(r.Context(), payload, r)

	m.next.ServeHTTP(w, r.WithContext(c))
}

//SetNext handler
func (m *ValidationMiddleware) SetNext(next http.Handler) {
	m.next = next
}
