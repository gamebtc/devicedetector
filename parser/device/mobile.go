package device

import (
	"path/filepath"
)

const ParserNameMobile = `mobile`
const FixtureFileMobile = `mobiles.yml`

func init() {
	RegDeviceParser(ParserNameMobile,
		func(dir string) DeviceParser {
			return NewMobile(filepath.Join(dir, FixtureFileMobile))
		})
}

func NewMobile(fileName string) *Mobile {
	m := &Mobile{}
	if err := m.Load(fileName); err != nil {
		return nil
	}
	return m
}

// Device parser for mobile detection
type Mobile struct {
	DeviceParserAbstract
}
