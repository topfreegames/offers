// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package errors

import "encoding/json"

//OfferTemplateError happens with invalid template
type OfferTemplateError struct {
	SourceError error
}

//NewOfferTemplateError constructor
func NewOfferTemplateError(err error) *OfferTemplateError {
	return &OfferTemplateError{
		SourceError: err,
	}
}

func (e *OfferTemplateError) Error() string {
	return e.SourceError.Error()
}

//Serialize returns the error serialized
func (e *OfferTemplateError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "OFF-003",
		"error":       "OfferTemplateError",
		"description": e.Error(),
	})

	return g
}
