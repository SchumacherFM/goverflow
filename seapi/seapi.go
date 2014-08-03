// seapi is short version for Stack Exchange API
package seapi

import (
	"net/http"
	"net/url"
	"time"
	"strings"
	"io/ioutil"
	"encoding/json"
	"errors"
)

const (
	API_TIMEOUT = 60
)

type Seapi struct {
	UseSsl        bool
	host          string
	version       string
	site          string
	transport     http.RoundTripper
	currentMethod []string
	currentParams url.Values
}

func NewSeapi(site string) (seapi *Seapi) {
	seapi = &Seapi{
		UseSsl: true,
		host:"://api.stackexchange.com",
		version: "2.2",
		site: site,
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
func (s *Seapi) SetMethod(method []string) *Seapi {
	s.currentMethod = method
	return s
}

func (s *Seapi) getQueryUrl() (theUrl string) {
	theUrl = "http"
	if true == s.UseSsl {
		theUrl = "https"
	}
	theUrl = theUrl+s.host+"/"+strings.Join(s.currentMethod, "/")

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
	for key, _ := range s.currentParams {
		delete(s.currentParams, key)
	}
}

func (s *Seapi) Query(collection interface{}) (error error) {

	client := &http.Client{
		Transport: s.getTransport(),
		Timeout: time.Second * API_TIMEOUT,
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
