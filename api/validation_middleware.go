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

	"github.com/asaskevich/govalidator"
	"github.com/topfreegames/offers/errors"
	"github.com/topfreegames/offers/models"
)

//ValidationMiddleware adds the version to the request
type ValidationMiddleware struct {
	GetPayload func() interface{}
	next       http.Handler
}

const payloadString = "payload"

func newContextWithPayload(payload interface{}, ctx context.Context, r *http.Request) context.Context {
	c := context.WithValue(ctx, payloadString, payload)
	return c
}

func gameFromCtx(ctx context.Context) *models.Game {
	game := ctx.Value(payloadString)
	if game == nil {
		return nil
	}
	return game.(*models.Game)
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
		WriteBytes(w, http.StatusBadRequest, vErr.Serialize())
		return
	}

	c := newContextWithPayload(payload, r.Context(), r)

	m.next.ServeHTTP(w, r.WithContext(c))
}

//SetNext handler
func (m *ValidationMiddleware) SetNext(next http.Handler) {
	m.next = next
}
