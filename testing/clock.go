// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package testing

import "time"

//MockClock allows for better control of time
type MockClock struct {
	CurrentTime int64
}

//GetTime returns mock time
func (m MockClock) GetTime() time.Time {
	return time.Unix(m.CurrentTime, 0)
}
