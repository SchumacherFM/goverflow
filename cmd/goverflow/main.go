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
	"github.com/codegangsta/cli"
	. "github.com/SchumacherFM/goverflow"
	"os"

)

func main() {
	app := cli.NewApp()
	app.Name = "goverflow"
	app.Version = "0.0.1"
	app.Usage = "Searches the stackexchange API and tweets new questions. App runs in the background or daemon."
	app.Action = showHelp
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "seconds,s",
			Value: 10,
			Usage: "Sleep duration in Seoncds, recommended: (3600*24)/300; quota is 300 queries / day",
		},
		cli.IntFlag{
			Name:  "logLevel,ll",
			Value: 0,
			Usage: "0 Debug, 1 Info, 2 Notice -> 7 Emergency",
		},
		cli.StringFlag{
			Name:  "logFile,lf",
			Value: "",
			Usage: "Log to file or if empty to os.Stderr",
		},
		cli.StringFlag{
			Name:  "configFile,c",
			Value: "config.json",
			Usage: "The JSON config file",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "Run the gin proxy in the current working directory",
			Action:    mainAction,
		},
	}

	app.Run(os.Args)
}

func showHelp(c *cli.Context) {
	cli.ShowAppHelp(c)
}

func mainAction(c *cli.Context) {

	a := NewGoverflowApp()
	a.SetInterval(c.GlobalInt("seconds"))
	a.SetLogFile(c.GlobalString("logFile"), c.GlobalInt("logLevel"))
	a.SetConfigFileName(c.GlobalString("configFile"))
	a.Goverflow()
}
