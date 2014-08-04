package poster

import (
	"encoding/json"
	"github.com/SchumacherFM/goverflow/seapi"
	"log"
	"net/url"
	"os"
)

type poster struct {
	logger *log.Logger
	Config struct {
		Host         string
		ApiVersion   string
		SearchParams string
	}
	so *seapi.Seapi
}

// Routineposter runs in a go routine
func (p *poster) RoutinePoster() {

	p.logger.Printf("%#v\n\n",p)


	p.logger.Println("Tick ...")
}

func NewPoster(fileName *string) *poster {
	p := &poster{
		so : seapi.NewSeapi(),
	}

	parseJsonConfig(p, fileName)

	p.so.Host = p.Config.Host
	p.so.Version = p.Config.ApiVersion

	parsed, err := url.Parse("http://dummy.com/?" + p.Config.SearchParams)
	if nil != err {
		panic(err)
	}
	p.so.SetParams(parsed.Query())
	p.so.SetMethod([]string{"search"})
	return p
}

func (p *poster) SetLogger(lg *log.Logger) {
	p.logger = lg
}

// parseJsonConfig parses the json file ;-)
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
