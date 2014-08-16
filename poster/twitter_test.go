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
	"github.com/SchumacherFM/goverflow/seapi"
	"testing"

	"fmt"
)

var (
	twitter = NewTwitter()
)

func init() {
	//	for logger see poster_test.go
	twitter.InitClient(logger, "testData/tweet.tpl.txt")
}

func TestGetTweet(t *testing.T) {
	logTestCollector = nil
	sr := &seapi.SearchResult{
		Link:  "http://gotest.com",
		Title: "Golang test tweet",
	}

	tweet, err := twitter.getTweet(sr)

	fmt.Printf("%#v %s", tweet, err)
	fmt.Print(logTestCollector)
}
