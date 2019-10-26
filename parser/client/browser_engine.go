package client

import (
	"path/filepath"

	. "github.com/gamebtc/devicedetector/parser"
)

// Known browser engines mapped to their internal short codes
var availableEngines = []string{
	`WebKit`,
	`Blink`,
	`Trident`,
	`Text-based`,
	`Dillo`,
	`iCab`,
	`Elektra`,
	`Presto`,
	`Gecko`,
	`KHTML`,
	`NetFront`,
	`Edge`,
	`NetSurf`,
}

const ParserNameBrowserEngine = `browserengine`
const FixtureFileBrowserEngine = `browser_engine.yml`

func init() {
	RegClientParser(ParserNameBrowserEngine,
		func(dir string) ClientParser {
			return NewBrowserEngine(filepath.Join(dir, FixtureFileBrowserEngine))
		})
}

func NewBrowserEngine(fileName string) *BrowserEngine {
	c := &BrowserEngine{}
	c.ParserName = ParserNameBrowserEngine
	if err := c.Load(fileName); err != nil {
		return nil
	}
	return c
}

type BrowserEngine struct {
	ClientParserAbstract
}

func (d *BrowserEngine) Parse(ua string) *ClientMatchResult {
	for _, regex := range d.Regexes {
		matches := regex.MatchUserAgent(ua)
		if len(matches) > 0 {
			name := BuildByMatch(regex.Name, matches)
			for _, v := range availableEngines {
				if StringEqualIgnoreCase(name, v) {
					return &ClientMatchResult{
						Type: ParserNameBrowserEngine,
						Name: v,
					}
				}
			}
		}
	}
	return nil
}
