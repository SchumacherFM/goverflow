/*
	Copyright (C) 2014  Cyrill AT Schumacher dot fm

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

    Contribute @ https://github.com/SchumacherFM/goverflow
*/

// seapi is short Version for Stack Exchange API
package seapi

import (
	"encoding/json"
	"errors"
	httpclient "github.com/SchumacherFM/goverflow/go-httpclient"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	API_TIMEOUT             = 60
	CONNECT_TIMEOUT         = 10
	REQUEST_TIMEOUT         = 30
	RESPONSE_HEADER_TIMEOUT = 20
)

type SeapiInterface interface {
	SetTransport(t http.RoundTripper) *Seapi
	getTransport() http.RoundTripper
	AddParam(key, value string) *Seapi
	SetParams(values url.Values) *Seapi
	SetMethod(method []string) *Seapi
	getQueryUrl() (theUrl string)
	resetQuery()
	Query(collection interface{}) (error error)
}

type Seapi struct {
	Host          string
	Version       string
	transport     http.RoundTripper
	currentMethod []string
	currentParams url.Values
}

func NewSeapi() (seapi *Seapi) {
	seapi = &Seapi{
		Host:          "https://api.stackexchange.com",
		Version:       "2.2",
		currentParams: make(url.Values),
	}

	return
}

func (s *Seapi) SetTransport(t http.RoundTripper) *Seapi {
	s.transport = t
	return s
}

func (s *Seapi) getTransport() http.RoundTripper {
	if nil == s.transport {
		s.transport = &httpclient.Transport{
			ConnectTimeout:        CONNECT_TIMEOUT * time.Second,
			RequestTimeout:        REQUEST_TIMEOUT * time.Second,
			ResponseHeaderTimeout: RESPONSE_HEADER_TIMEOUT * time.Second,
		}
	}

	return s.transport
}

func (s *Seapi) AddParam(key, value string) *Seapi {
	s.currentParams.Add(key, value)
	return s
}
func (s *Seapi) SetParam(key, value string) *Seapi {
	s.currentParams.Set(key, value)
	return s
}
func (s *Seapi) SetParams(values url.Values) *Seapi {
	s.currentParams = values
	return s
}
func (s *Seapi) SetMethod(method []string) *Seapi {
	s.currentMethod = method
	return s
}

func (s *Seapi) getQueryUrl() (theUrl string) {
	theUrl = s.Host + "/" + strings.Join(s.currentMethod, "/")
	if len(s.currentParams) > 0 {
		theUrl = theUrl + "?" + s.currentParams.Encode()
	}

	return
}

func (s *Seapi) ResetQuery() {
	if nil != s.currentMethod {
		// s.currentMethod[:0]
		s.currentMethod = nil
	}
	for key := range s.currentParams {
		delete(s.currentParams, key)
	}
}

func (s *Seapi) Query(collection interface{}) (qryUrl string, error error) {

	client := &http.Client{
		Transport: s.getTransport(),
		Timeout:   time.Second * API_TIMEOUT,
	}
	qryUrl = s.getQueryUrl()
	response, error := client.Get(qryUrl)
	defer response.Body.Close()
	if nil != error {
		return
	}
	if 200 != response.StatusCode {
		return "", errors.New(response.Status + " @ " + qryUrl)
	}

	apiContent, error := ioutil.ReadAll(response.Body)

	if error != nil {
		return
	}
	error = json.Unmarshal(apiContent, collection)
	if error != nil {
		return
	}
	return
}
