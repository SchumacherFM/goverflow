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
	"errors"
	"fmt"
	log "github.com/SchumacherFM/goverflow/go-log"
	"github.com/SchumacherFM/goverflow/seapi"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"net/http"
	"net/url"
	"strings"
	"bytes"
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
}

func (t *Twitter) InitClient(logger *log.Logger) {
	t.logger = logger

	config := &oauth1a.ClientConfig{
		ConsumerKey:    t.ConsumerKey,
		ConsumerSecret: t.ConsumerSecret,
	}
	user := oauth1a.NewAuthorizedConfig(t.AccessToken, t.AccessTokenSecret)
	t.client = twittergo.NewClient(config, user)
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

	data.Set("status", getTweet(sr))
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

func getTweet(sr *seapi.SearchResult) string {
	tweetSuffix := "\n#golang" // @todo make configurable
	tweet := bytes.NewBufferString(sr.Title)
	maxLen := TWEET_LENGTH - TCO_LENGTH - len(tweetSuffix) - 2 // 2 is whitespaces
	if tweet.Len() > maxLen {
		tweet.Truncate(maxLen)
	}
	tweet.WriteString("\n")
	tweet.WriteString(sr.Link)
	tweet.WriteString(tweetSuffix)
	return tweet.String()
}
