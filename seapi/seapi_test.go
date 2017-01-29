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
// seapi is short version for Stack Exchange API
package seapi

import (
	"strings"
	"testing"

	"github.com/SchumacherFM/goverflow/testHelper"
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

	roundTripper, httpTS := testHelper.GetMockServerForPath("/search", "test_search1", t)
	defer httpTS.Close()
	seapi.SetTransport(roundTripper)

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
