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
package main

import (
	. "github.com/SchumacherFM/goverflow"
	"flag"

)

func main() {
	inputDuration := flag.Int("seconds", 10, "Sleep duration in Seoncds, recommended: (3600*24)/300; quota is 300 queries")
	logLevel := flag.Int("logLevel", 0, "0 Debug, 1 Info, 2 Notice -> 7 Emergency")
	logFile := flag.String("logFile", "", "Log to file or if empty to os.Stderr")
	configFile := flag.String("configFile", "config.json", "Config file")
	flag.Parse()

	a := NewGoverflowApp()

	a.SetInterval(inputDuration)
	a.SetLogFile(logFile, logLevel)
	a.SetConfigFileName(configFile)
	a.Goverflow()
}
