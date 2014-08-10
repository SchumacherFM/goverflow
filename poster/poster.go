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
	"encoding/json"
	log "github.com/SchumacherFM/goverflow/go-log"
	"github.com/SchumacherFM/goverflow/seapi"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"
)

type poster struct {
	logger *log.Logger
	Config struct {
		Host         string
		ApiVersion   string
		SearchParams string
	}
	timeLastRun      int64
	quotaRemaining   int
	so             *seapi.Seapi
	gfdb             GFDB
}

func NewPoster(fileName *string) *poster {

	var err error
	db := new(GFDB)
	err = db.InitDb()
	if nil != err {
		panic("failed to create memDB: " + err.Error())
	}

	p := &poster{
		so: seapi.NewSeapi(),
		gfdb: *db,
	}

	parseJsonConfig(p, fileName)

	p.so.Host = p.Config.Host
	p.so.Version = p.Config.ApiVersion
	p.setTimeLastRun()

	parsed, err := url.Parse("http://dummy.com/?" + p.Config.SearchParams)
	if nil != err {
		panic(err)
	}
	p.so.SetParams(parsed.Query())
	p.so.SetMethod([]string{"search"})
	return p
}

// Routineposter runs in a go routine
func (p *poster) RoutinePoster() error {
	defer p.setTimeLastRun()
	soSearchResultCollection := p.routineGetSearchCollection()
	if nil == soSearchResultCollection {
		return nil // no further processing
	}

	// sort map of soSearchResultCollection by lowest question id to highest
	var keys []int
	for k := range soSearchResultCollection {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		sr := soSearchResultCollection[k]
		// now post to twitter and set value in DB
		p.logger.Debug("RC: %v", sr)
	}

	p.logger.Debug("Tick ...\n")
	return nil
}

// routineGetCollection runs within a goroutine
func (p *poster) routineGetSearchCollection() map[int]seapi.SearchResult {
	soSearchResultCollection := &seapi.SearchResultCollection{}

	// change fromDate to the current timeStamp of this run
	p.so.SetParam("fromdate", strconv.FormatInt(p.timeLastRun, 10))

	queryUrl, err := p.so.Query(soSearchResultCollection)
	p.logger.Debug("Query URL: %s", queryUrl)
	if nil != err {
		p.logger.Emergency("L36: %s", err.Error())
		return nil // no further processing in this routine
	}

	if 0 == len(soSearchResultCollection.Items) {
		p.logger.Debug("No new questions posted since %s", p.getTimeLastRunRFC1123Z())
		return nil
	}

	if 0 == soSearchResultCollection.Quota_remaining {
		p.logger.Debug("Over quota :-( %s", p.getTimeLastRunRFC1123Z())
		return nil
	} else {
		p.logger.Debug("Quota remaining: %d", soSearchResultCollection.Quota_remaining)
	}

	// now calculate the difference and return only the new items
	var newItems = make(map[int]seapi.SearchResult)
	for _, searchResult := range soSearchResultCollection.Items {
		storedResult, err := p.gfdb.FindByQuestionId(searchResult.Question_id)
		if nil != err {
			p.logger.Error("FindByQuestionId: %s", err)
		}
		if nil != storedResult { // already posted
			continue
		}
		newItems[searchResult.Question_id] = searchResult
	}

	return newItems
}

func (p *poster) setTimeLastRun() {
	p.timeLastRun = time.Now().Unix()-3600*12 // last part just for testing
}

func (p *poster) getTimeLastRunRFC1123Z() string {
	return time.Unix(p.timeLastRun, 0).Format(time.RFC1123Z)
}

func (p *poster) SetLogger(lg *log.Logger) {
	p.logger = lg
}

// parseJsonConfig parses the json file ;-) no logger available
func parseJsonConfig(p *poster, fileName *string) {
	file, err := os.Open(*fileName)

	if nil != err {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&p.Config)
	if nil != err {
		panic(err)
	}
}
