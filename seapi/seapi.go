// seapi is short Version for Stack Exchange API
package seapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	API_TIMEOUT = 60
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
		s.transport = http.DefaultTransport
	}
	return s.transport
}

func (s *Seapi) AddParam(key, value string) *Seapi {
	s.currentParams.Add(key, value)
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
	theUrl = s.Host+"/"+strings.Join(s.currentMethod, "/")
	if len(s.currentParams) > 0 {
		theUrl = theUrl+"?"+s.currentParams.Encode()
	}

	return
}

func (s *Seapi) resetQuery() {
	if nil != s.currentMethod {
		// s.currentMethod[:0]
		s.currentMethod = nil
	}
	for key := range s.currentParams {
		delete(s.currentParams, key)
	}
}

func (s *Seapi) Query(collection interface{}) (error error) {

	client := &http.Client{
		Transport: s.getTransport(),
		Timeout:   time.Second * API_TIMEOUT,
	}
	qryUrl := s.getQueryUrl()
	response, error := client.Get(qryUrl)
	s.resetQuery()
	if nil != error {
		return error
	}
	if 200 != response.StatusCode {
		return errors.New(response.Status + " @ " + qryUrl)
	}

	defer response.Body.Close()
	apiContent, error := ioutil.ReadAll(response.Body)

	if error != nil {
		return error
	}
	error = json.Unmarshal(apiContent, collection)
	if error != nil {
		return error
	}
	return nil
}
