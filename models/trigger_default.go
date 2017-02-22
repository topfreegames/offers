// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

//DefaultTrigger implements interface Trigger
type DefaultTrigger struct{}

//IsTriggered returns the current time
func (dt DefaultTrigger) IsTriggered(times map[string]interface{}, user map[string]interface{}) bool {
  return true
}
