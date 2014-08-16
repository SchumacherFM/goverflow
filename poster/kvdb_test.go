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
	"github.com/kurrik/twittergo"
	"testing"
)

var (
	dbTest  = NewGFDB()
	tweetId = 4711 // echt koelnisch wasser
	aTweet  = make(twittergo.Tweet)
)

func init() {
	aTweet["id_str"] = (interface{})("123456789")
	aTweet["text"] = (interface{})("Es chunnt scho guet") // Swiss German ;-)
}

func TestSaveTweet(t *testing.T) {

	err := dbTest.SaveTweet(tweetId, &aTweet)
	if nil != err {
		t.Error("cannot save tweet", err)
	}

}
func TestFindByQuestionId(t *testing.T) {

	result, err := dbTest.FindByQuestionId(tweetId)
	if nil != err {
		t.Error("cannot find tweet", err)
	}

	if "123456789;Es chunnt scho guet" != string(result) {
		t.Error("Cannot check that '123456789;Es chunnt scho guet' is in ", string(result))
	}

}
