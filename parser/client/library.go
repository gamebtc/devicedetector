package client

import (
	"path/filepath"
)

const ParserNameLibrary = `library`
const FixtureFileLibrary = `libraries.yml`

func init() {
	RegClientParser(ParserNameLibrary,
		func(dir string) ClientParser {
			return NewLibrary(filepath.Join(dir, FixtureFileLibrary))
		})
}

func NewLibrary(fileName string) *Library {
	c := &Library{}
	c.ParserName = ParserNameLibrary
	if err := c.Load(fileName); err != nil {
		return nil
	}
	return c
}

// Client parser for tool & software detection
type Library struct {
	ClientParserAbstract
}
