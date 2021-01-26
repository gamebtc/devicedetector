package device

import (
	"path/filepath"
	"testing"

	. "github.com/muxinc/devicedetector/parser"
	"gotest.tools/assert"
)

func TestConsoleParse(t *testing.T) {
	ps := NewConsole(filepath.Join(dir, FixtureFileConsole))
	var list []*DeviceFixture
	err := ReadYamlFile(`fixtures/console.yml`, &list)
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

