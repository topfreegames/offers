// offers api
// https://github.com/topfreegames/offers
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import "net/http"

//Write to the response and with the status code
func Write(w http.ResponseWriter, status int, text string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	WriteBytes(w, status, []byte(text))
}

//WriteBytes to the response and with the status code
func WriteBytes(w http.ResponseWriter, status int, text []byte) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(text)
}
