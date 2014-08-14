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
	"fmt"
	"github.com/cznic/kv"
	"github.com/kurrik/twittergo"
)

const (
	DB_TWEET_SEP = ";"
)

type GFDB struct {
	*kv.DB
}

func NewGFDB() *GFDB {
	db := new(GFDB)

	err := db.initDb()
	if nil != err {
		panic("failed to create memDB: " + err.Error())
	}
	return db
}

func (db *GFDB) initDb() error {
	var err error
	kvOpt := &kv.Options{}
	db.DB, err = kv.CreateMem(kvOpt)
	return err
}

func (db *GFDB) FindByQuestionId(id int) ([]byte, error) {
	return db.Get(nil, db.makeQuestionKey(id))
}

func (db *GFDB) SaveTweet(id int, tweet *twittergo.Tweet) error {

	var value bytes.Buffer
	value.WriteString(tweet.IdStr())
	value.WriteString(DB_TWEET_SEP)
	value.WriteString(tweet.Text())

	return db.Set(db.makeQuestionKey(id), value.Bytes())
}

func (db *GFDB) makeQuestionKey(id int) []byte {
	return []byte(fmt.Sprintf("q_%d", id))
}
