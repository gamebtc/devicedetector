package client

import (
	"path/filepath"
	"testing"

	. "github.com/gamebtc/devicedetector/parser"
	"gotest.tools/assert"
)

func TestLibraryParse(t *testing.T) {
	var ps= NewLibrary(filepath.Join(dir, FixtureFileLibrary))
	var list []*ClientFixture
	err := ReadYamlFile(`fixtures/library.yml`, &list)
	if err != nil {
		t.Error(err)
	}

	for _, item := range list {
		ua := item.UserAgent
		r := ps.Parse(ua)
		assert.DeepEqual(t, item.ClientMatchResult, r)
	}
}
