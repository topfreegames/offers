// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Models <backend@tfgco.com>

package errors

import (
	"encoding/json"
	"fmt"
)

//ModelNotFoundError happens when game could not be found with specified arguments
type ModelNotFoundError struct {
	Model   string
	Filters map[string]interface{}
}

//NewModelNotFoundError ctor
func NewModelNotFoundError(model string, filters map[string]interface{}) *ModelNotFoundError {
	return &ModelNotFoundError{
		Model:   model,
		Filters: filters,
	}
}

func (e *ModelNotFoundError) Error() string {
	return fmt.Sprintf("%s was not found with specified filters.", e.Model)
}

//Serialize returns the error serialized
func (e *ModelNotFoundError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "OFF-001",
		"error":       fmt.Sprintf("%sNotFoundError", e.Model),
		"description": e.Error(),
		"filters":     e.Filters,
	})

	return g
}

//InvalidModelError happens when a model could not be saved/updated/deleted
type InvalidModelError struct {
	Model   string
	Message string
}

//NewInvalidModelError ctor
func NewInvalidModelError(model, message string) *InvalidModelError {
	return &InvalidModelError{
		Model:   model,
		Message: message,
	}
}

func (e *InvalidModelError) Error() string {
	return fmt.Sprintf("%s could not be saved due to: %s", e.Model, e.Message)
}

//Serialize returns the error serialized
func (e *InvalidModelError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "OFF-003",
		"error":       fmt.Sprintf("Invalid%sError", e.Model),
		"description": e.Error(),
	})

	return g
}
