package device

import (
	"path/filepath"
	"testing"

	. "github.com/gamebtc/devicedetector/parser"
	"gotest.tools/assert"
)

func TestCarParse(t *testing.T) {
	ps := NewCar(filepath.Join(dir, FixtureFileCar))
	var list []*DeviceFixture
	err := ReadYamlFile(`fixtures/car_browser.yml`, &list)
	if err != nil {
		t.Error(err)
	}

	for _, item := range list {
		ua := item.UserAgent
		r := ps.Parse(ua)
		test := item.GetDeviceMatchResult()
		assert.DeepEqual(t, test, r)
	}
}
