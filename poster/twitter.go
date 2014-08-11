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
	log "github.com/SchumacherFM/goverflow/go-log"
	"github.com/SchumacherFM/goverflow/seapi"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"
)

const (
	TWEET_LENGTH = 140
	TCO_LENGTH   = 20
)

type Twitter struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	client            *twittergo.Client
	logger            *log.Logger
	tweetTpl          *template.Template
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

// TweetQuestion posts a tweet to twitter
// @see https://dev.twitter.com/docs/api/1/post/statuses/update
// returns error
func (t *Twitter) TweetQuestion(sr *seapi.SearchResult) error {
	var (
		err   error
		req   *http.Request
		resp  *twittergo.APIResponse
		tweet *twittergo.Tweet
	)

	data := url.Values{}
	data.Set("status", t.getTweet(sr))
	body := strings.NewReader(data.Encode())
	req, err = http.NewRequest("POST", "/1.1/statuses/update.json", body)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not parse request: %v\n", err))

	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = t.client.SendRequest(req)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not send request: %v\n", err))
	}
	tweet = &twittergo.Tweet{}
	err = resp.Parse(tweet)
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
		return err
	}

	fmt.Printf("ID:                   %v\n", tweet.Id())
	fmt.Printf("Tweet:                %v\n", tweet.Text())
	fmt.Printf("User:                 %v\n", tweet.User().Name())

	if true == resp.HasRateLimit() {
		fmt.Printf("Rate limit:           %v\n", resp.RateLimit())
		fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
		fmt.Printf("Rate limit reset:     %v\n", resp.RateLimitReset())
	} else {
		fmt.Printf("Could not parse rate limit from response.\n")
	}

	return nil
}

func (t *Twitter) getTweet(sr *seapi.SearchResult) string {
	var theTweet bytes.Buffer

	err := t.tweetTpl.Execute(&theTweet, sr)

	if nil != err {
		t.logger.Emergency("Error template %s", err)
	}
	// todo check for length and convert e.g. &quot; into "
	return theTweet.String()
	//	maxLen := TWEET_LENGTH - TCO_LENGTH - len(tweetSuffix) - 2 // 2 is whitespaces

}
