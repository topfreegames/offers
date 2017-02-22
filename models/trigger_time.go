// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

//TimeTrigger implements interface Trigger
type TimeTrigger struct{}

//Times holds from and to in UnixTimestamp
type Times struct {
  from int64
  to int64
}

//IsTriggered returns the current time
func (tt *TimeTrigger) IsTriggered(times interface{}, now interface{}) bool {
  t := times.(Times)
  n := now.(int64)

  return t.from <= n && n <= t.to
}
