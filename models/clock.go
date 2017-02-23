// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package models

import "time"

//RealClock returns the actual time from the OS
type RealClock struct{}

//GetTime returns the current time
func (r RealClock) GetTime() time.Time {
	return time.Now()
}
