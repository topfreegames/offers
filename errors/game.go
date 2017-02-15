// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package errors

import "encoding/json"

//SerializableError means that an error can be transformed to JSON
type SerializableError interface {
	Serialize() []byte
}

//GameNotFoundError happens when game could not be found with specified arguments
type GameNotFoundError struct {
	Filters map[string]interface{}
}

//NewGameNotFoundError ctor
func NewGameNotFoundError(filters map[string]interface{}) *GameNotFoundError {
	return &GameNotFoundError{
		Filters: filters,
	}
}

func (e *GameNotFoundError) Error() string {
	return "Game was not found with specified filters."
}

//Serialize returns the error serialized
func (e *GameNotFoundError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "OFF-01",
		"error":       "GameNotFoundError",
		"description": "Game was not found with specified filters.",
		"filters":     e.Filters,
	})

	return g
}
