package device

import (
	"path/filepath"
)

const ParserNameConsole = `console`
const FixtureFileConsole = `consoles.yml`

func init() {
	RegDeviceParser(ParserNameConsole,
		func(dir string) DeviceParser {
			return NewConsole(filepath.Join(dir, FixtureFileConsole))
		})
}

func NewConsole(fileName string) *Console {
	c := &Console{}
	if err := c.Load(fileName); err != nil {
		return nil
	}
	return c
}

// Device parser for console detection
type Console struct {
	DeviceParserAbstract
}

func (c *Console) Parse(ua string) *DeviceMatchResult {
	if !c.PreMatch(ua) {
		return nil
	}
	return c.DeviceParserAbstract.Parse(ua)
}
