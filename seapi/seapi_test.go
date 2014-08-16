// seapi is short version for Stack Exchange API
package seapi

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var (
	seapi = NewSeapi()
)

func TestGetQueryUrlEmptyHttps(t *testing.T) {
	actual := seapi.getQueryUrl()
	if expected := "https://api.stackexchange.com/"; expected != actual {
		t.Error("Expected:", expected, " but got: ", actual)
	}
}

func TestQueryPrepare1(t *testing.T) {
	seapi.AddParam("order", "desc")
	seapi.AddParam("sort", "creation")
	seapi.SetMethod([]string{"answers", "2;34;43", "comments"})

	actual := seapi.getQueryUrl()
	expected := "https://api.stackexchange.com/answers/2;34;43/comments?order=desc&sort=creation"
	if expected != actual {
		t.Error("Expected:", expected, " but got: ", actual)
	}
	seapi.ResetQuery()
}

func TestQueryPrepare2(t *testing.T) {
	seapi.AddParam("order", "asc")
	seapi.AddParam("sort", "com ment")
	seapi.SetMethod([]string{"questions", "2;314;43"})

	actual := seapi.getQueryUrl()
	expected := "https://api.stackexchange.com/questions/2;314;43?order=asc&sort=com+ment"
	if expected != actual {
		t.Error("Expected:", expected, " but got: ", actual)
	}
	seapi.ResetQuery()
}

func TestResetQuery(t *testing.T) {
	seapi.AddParam("order", "asc")
	seapi.AddParam("sort", "com ment")
	seapi.SetMethod([]string{"questions", "2,314,43", "comments"})
	len1 := len(seapi.currentParams)
	if len1 != 2 {
		t.Error("CurrentParams: Expected 3 params but got ", len1)
	}
	len1a := len(seapi.currentMethod)
	if len1a != 3 {
		t.Error("CurrentMethod Expected 3 params but got ", len1a)
	}
	seapi.ResetQuery()
	len2 := len(seapi.currentParams)
	if len2 != 0 {
		t.Error("CurrentParams: Expected 0 params but got ", len2)
	}
	len2a := len(seapi.currentMethod)
	if len2a != 0 {
		t.Error("CurrentMethod Expected 0 params but got ", len2a)
	}
}

func TestQuery1(t *testing.T) {
	seapi.Host = "http://api.stackexchange.com"
	seapi.AddParam("order", "desc")
	seapi.AddParam("sort", "creation")
	seapi.SetMethod([]string{"search"})

	httpTS := getMockServerForPath("/search", "test_search1", t)
	defer httpTS.Close()

	searchResult := &SearchResultCollection{}
	qryUrl, qryErr := seapi.Query(searchResult)

	if nil != qryErr {
		t.Fatal(qryErr.Error())
	} else {
		if len(searchResult.Items) == 0 {
			t.Error("SearchResultCollection.Items is zero :-(")
		}
	}

	if false == strings.Contains(qryUrl, ".com/search") {
		t.Error("Cannot find .com/search in qryUrl")
	}

	if 300 != searchResult.Quota_max {
		t.Error("Cannot parse Quota Max")
	}

}

// getMockServerForPath creates a test server and sets a the transport proxy server to the
// test server URL. So any request to any URL will result in an answer from the proxy
// testData will be written as the response
func getMockServerForPath(path string, testDataJsonFile string, t *testing.T) *httptest.Server {

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
	seapi.SetTransport(roundTripper)
	return ts
}
