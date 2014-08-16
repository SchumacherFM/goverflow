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
package testHelper

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// getMockServerForPath creates a test server and sets a the transport proxy server to the
// test server URL. So any request to any URL will result in an answer from the proxy
// testData will be written as the response
func GetMockServerForPath(path string, testDataJsonFile string, t *testing.T) (http.RoundTripper, *httptest.Server) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			t.Error("Path doesn't match. Expected: ", path, " Actual: ", r.URL.Path)
			http.Error(w, "Path doesn't match", http.StatusInternalServerError)
		} else {
			testData, err := ioutil.ReadFile("./testData/" + testDataJsonFile + ".json")
			if nil != err {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Write(testData)
			}
		}
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	roundTripper := &http.Transport{
		Proxy: func(*http.Request) (*url.URL, error) {
			//change the host to use the test server http://127.0.0.1:XXXX; could be any port
			return url.Parse(ts.URL)
		},
	}

	return roundTripper, ts
}
