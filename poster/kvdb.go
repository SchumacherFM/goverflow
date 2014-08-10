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
	"fmt"
	"github.com/SchumacherFM/goverflow/kv"
)

type GFDB struct {
	*kv.DB
}

func (db *GFDB) InitDb() error {
	var err error
	kvOpt := &kv.Options{}
	db.DB, err = kv.CreateMem(kvOpt)
	return err
}

func (db *GFDB) FindByQuestionId(id int) ([]byte, error) {
	return db.Get(nil, db.makeQuestionKey(id))
}

func (db *GFDB) makeQuestionKey(id int) []byte {
	return []byte(fmt.Sprintf("q_%d", id))
}
