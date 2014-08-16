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

package seapi

type Owner struct {
	Reputation    int
	User_id       int
	User_type     string //one of unregistered, registered, moderator, or does_not_exist
	Accept_rate   int
	Profile_image string
	Display_name  string
	Link          string
}

type SearchResult struct {
	Tags               []string
	Owner              Owner
	Is_answered        bool
	View_count         int
	Answer_count       int
	Score              int
	Last_activity_date int64
	Creation_date      int64
	Question_id        int
	Link               string
	Title              string
}

type SearchResultCollection struct {
	Items           []SearchResult
	Has_more        bool
	Quota_max       int
	Quota_remaining int

	Error_id      int
	Error_name    string
	Error_message string
}
