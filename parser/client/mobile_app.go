package client

import (
	"path/filepath"
)

const ParserNameMobileApp = `mobile app`
const FixtureFileMobileApp = `mobile_apps.yml`

func init() {
	RegClientParser(ParserNameMobileApp,
		func(dir string) ClientParser {
			return NewMobileApp(filepath.Join(dir, FixtureFileMobileApp))
		})
}

func NewMobileApp(fileName string) *MobileApp {
	c := &MobileApp{}
	c.ParserName = ParserNameMobileApp
	if err := c.Load(fileName); err != nil {
		return nil
	}
	return c
}

// Client parser for mobile app detection
type MobileApp struct {
	ClientParserAbstract
}
