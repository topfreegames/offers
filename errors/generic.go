// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package errors

import "encoding/json"

//GenericError happens when an unidentified error occurs
type GenericError struct {
	Message     string
	SourceError error
}

//NewGenericError ctor
func NewGenericError(message string, err error) *GenericError {
	return &GenericError{
		Message:     message,
		SourceError: err,
	}
}

func (e *GenericError) Error() string {
	return e.SourceError.Error()
}

//Serialize returns the error serialized
func (e *GenericError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "OFF-004",
		"error":       e.Message,
		"description": e.SourceError.Error(),
	})

	return g
}
