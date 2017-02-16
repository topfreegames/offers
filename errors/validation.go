// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package errors

import "encoding/json"

//ValidationFailedError happens when validation of a struct fails
type ValidationFailedError struct {
	SourceError error
}

//NewValidationFailedError ctor
func NewValidationFailedError(err error) *ValidationFailedError {
	return &ValidationFailedError{
		SourceError: err,
	}
}

func (e *ValidationFailedError) Error() string {
	return e.SourceError.Error()
}

//Serialize returns the error serialized
func (e *ValidationFailedError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "OFF-002",
		"error":       "ValidationFailedError",
		"description": e.SourceError.Error(),
	})

	return g
}
