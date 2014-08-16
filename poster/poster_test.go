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

package poster

import (
	"testing"
	"github.com/SchumacherFM/goverflow/testHelper"
	log "github.com/segmentio/go-log"
	"github.com/kurrik/twittergo"
	"github.com/SchumacherFM/goverflow/seapi"
	"net/url"
	"strings"
)

var (
	testPoster *poster
	fileName = "./testData/config.json"
	logger *log.Logger
	lm       = &LoggerMock{}
	logTestCollector []string
)

type LoggerMock struct {}

func (lm *LoggerMock) Write(p []byte) (n int, err error) {
	logTestCollector = append(logTestCollector, string(p))
	return 1, nil
}

type TwitterMock struct {

}

func (t *TwitterMock) GetTwitter() *Twitter {
	return NewTwitter()
}

func (t *TwitterMock) GetConfig() map[string]string {
	cm := make(map[string]string)
	return cm
}

func (t *TwitterMock) InitClient(logger *log.Logger, tweetTplFile string) {
}

func (t *TwitterMock) doRequest(data *url.Values) (*twittergo.Tweet, error) {
	var (
		err   error
		tweet *twittergo.Tweet
	)
	return tweet, err
}

func (t *TwitterMock) TweetQuestion(sr *seapi.SearchResult) (*twittergo.Tweet, error) {

	aTweet := make(twittergo.Tweet)

	aTweet["id_str"] = string(sr.Question_id)
	aTweet["text"] = sr.Title

	return &aTweet, nil
}

func (t *TwitterMock) getTweet(sr *seapi.SearchResult) (string, error) {
	return "", nil
}

func init() {
	logger = log.New(lm, 0, "Test")
	testPoster = NewPoster(&fileName, logger)
}

// tests the real poster instance
func TestNewPoster(t *testing.T) {

	if "http://api.stackexchange.com" != testPoster.Config.Host {
		t.Error("Not found http://api.stackexchange.com")
	}
	if "2.2" != testPoster.Config.ApiVersion {
		t.Error("Not found 2.2")
	}
	if "order=desc&sort=creation&tagged=go&site=stackoverflow" != testPoster.Config.SearchParams {
		t.Error("Not found order=desc&sort=creation&tagged=go&site=stackoverflow")
	}
	twitterConfig := testPoster.twitter.GetConfig()
	if "consumerKeyconsumerKey" != twitterConfig["ConsumerKey"] {
		t.Error("Not found consumerKeyconsumerKey")
	}
	if "consumerSecretconsumerSecret" != twitterConfig["ConsumerSecret"] {
		t.Error("Not found consumerSecretconsumerSecret")
	}
	if "accessTokenaccessToken" != twitterConfig["AccessToken"] {
		t.Error("Not found accessTokenaccessToken")
	}
	if "accessTokenSecretaccessTokenSecret" != twitterConfig["AccessTokenSecret"] {
		t.Error("Not found accessTokenSecretaccessTokenSecret")
	}

}

// uses the twitter mock
func TestRoutinePoster(t *testing.T) {
	testPoster.twitter = &TwitterMock{}
	testPoster.so.SetMethod([]string{"search"})
	roundTripper, httpTS := testHelper.GetMockServerForPath("/search", "test_search1", t)
	defer httpTS.Close()
	testPoster.so.SetTransport(roundTripper)
	testPoster.RoutinePoster()

	logData := strings.Join(logTestCollector, "\n")
	assertLogTestCollector(t, &logData, "Tick ...")
	assertLogTestCollector(t, &logData, "mgo: how to update a specific array in a document")
	assertLogTestCollector(t, &logData, "Quota remaining: 274")
}

func assertLogTestCollector(t *testing.T, logData *string, test string) {
	if false == strings.Contains(*logData, test) {
		t.Error("Cannot find:", test, "in", *logData)
	}
}
