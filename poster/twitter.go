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
	"bytes"
	"errors"
	"fmt"
	"github.com/SchumacherFM/goverflow/seapi"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	log "github.com/segmentio/go-log"
	"html"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
)

const (
	TWEET_LENGTH = 140
	TCO_LENGTH   = 24 // officially it's 20 chars ... but lets add 4 for back up
)

type TwitterInterface interface {
	GetTwitter() *Twitter
	GetConfig() map[string]string
	InitClient(logger *log.Logger, tweetTplFile string)
	doRequest(data *url.Values) (*twittergo.Tweet, error)
	TweetQuestion(sr *seapi.SearchResult) (*twittergo.Tweet, error)
}

type Twitter struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	client            *twittergo.Client
	logger            *log.Logger
	tweetTpl          *template.Template
}

func NewTwitter() *Twitter {
	return &Twitter{}
}

func (t *Twitter) GetTwitter() *Twitter {
	return t
}

func (t *Twitter) GetConfig() map[string]string {
	cm := make(map[string]string, 4)
	cm["ConsumerKey"] = t.ConsumerKey
	cm["ConsumerSecret"] = t.ConsumerSecret
	cm["AccessToken"] = t.AccessToken
	cm["AccessTokenSecret"] = t.AccessTokenSecret
	return cm
}

func (t *Twitter) InitClient(logger *log.Logger, tweetTplFile string) {
	var err error
	t.logger = logger

	config := &oauth1a.ClientConfig{
		ConsumerKey:    t.ConsumerKey,
		ConsumerSecret: t.ConsumerSecret,
	}
	user := oauth1a.NewAuthorizedConfig(t.AccessToken, t.AccessTokenSecret)
	t.client = twittergo.NewClient(config, user)

	t.tweetTpl, err = template.New(tweetTplFile).ParseFiles(tweetTplFile)
	if nil != err {
		t.logger.Emergency("Cannot open tweet template file %s", tweetTplFile)
		os.Exit(2)
	}
}
func (t *Twitter) doRequest(data *url.Values) (*twittergo.Tweet, error) {
	var (
		err   error
		req   *http.Request
		resp  *twittergo.APIResponse
		tweet *twittergo.Tweet
	)

	body := strings.NewReader(data.Encode())
	req, err = http.NewRequest("POST", "/1.1/statuses/update.json", body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not parse request: %v\n", err))

	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = t.client.SendRequest(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not send request: %v\n", err))
	}

	if true == resp.HasRateLimit() {
		t.logger.Warning("Rate limit:           %v\n", resp.RateLimit())
		t.logger.Warning("Rate limit remaining: %v\n", resp.RateLimitRemaining())
		t.logger.Warning("Rate limit reset:     %v\n", resp.RateLimitReset())
	} else {
		t.logger.Notice("Could not parse rate limit from response.")
	}

	tweet = &twittergo.Tweet{}
	err = resp.Parse(tweet)
	return tweet, err
}

// TweetQuestion posts a tweet to twitter
// @see https://dev.twitter.com/docs/api/1/post/statuses/update
// returns error
func (t *Twitter) TweetQuestion(sr *seapi.SearchResult) (*twittergo.Tweet, error) {
	var (
		err      error
		twString string
		tweet    *twittergo.Tweet
	)

	twString, err = t.getTweet(sr)
	if nil != err {
		return nil, err
	}
	data := &url.Values{}
	data.Set("status", twString)
	tweet, err = t.doRequest(data)

	if err != nil {
		if rle, ok := err.(twittergo.RateLimitError); ok {
			t.logger.Warning("Rate limited, reset at %v\n", rle.Reset)
		} else if errs, ok := err.(twittergo.Errors); ok {
			for i, val := range errs.Errors() {
				t.logger.Error("Error #%v - ", i+1)
				t.logger.Error("Code: %v ", val.Code())
				t.logger.Error("Msg: %v\n", val.Message())
			}
		} else {
			t.logger.Error("Problem parsing response: %v\n", err)
		}
		return nil, err
	}
	return tweet, nil
}

func (t *Twitter) getTweet(sr *seapi.SearchResult) (string, error) {

	if (len(sr.Title)+TCO_LENGTH) > TWEET_LENGTH { // just a simple check
		return "", errors.New("Tweet is too long ...")
	}

	var theTweet bytes.Buffer
	err := t.tweetTpl.Execute(&theTweet, sr)
	if nil != err {
		t.logger.Emergency("Error template %s", err)
	}
	// fix https://twitter.com/davecheney/status/499495512555266048
	return html.UnescapeString(theTweet.String()), nil

}
