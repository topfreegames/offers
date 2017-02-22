// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import (
	"time"
)

//TimeTrigger implements interface Trigger
type TimeTrigger struct{}

//Times holds from and to in UnixTimestamp
type Times struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

//IsTriggered returns the current time
func (tt TimeTrigger) IsTriggered(times interface{}, now interface{}) bool {
	t := times.(Times)
	n := now.(time.Time).Unix()

	return t.From <= n && n <= t.To
}
