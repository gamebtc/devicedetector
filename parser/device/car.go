package device

import (
	"path/filepath"
)

const ParserNameCar = `car browser`
const FixtureFileCar = `car_browsers.yml`

func init() {
	RegDeviceParser(ParserNameCar,
		func(dir string) DeviceParser {
			return NewCar(filepath.Join(dir, FixtureFileCar))
		})
}

func NewCar(fileName string) *Car {
	c := &Car{}
	if err := c.Load(fileName); err != nil {
		return nil
	}
	return c
}

// Device parser for car browser detection
type Car struct {
	DeviceParserAbstract
}

func (c *Car) Parse(ua string) *DeviceMatchResult {
	if !c.PreMatch(ua) {
		return nil
	}
	return c.DeviceParserAbstract.Parse(ua)
}
