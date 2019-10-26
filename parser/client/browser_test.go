package client

import (
	"path/filepath"
	"testing"

	. "github.com/gamebtc/devicedetector/parser"
	"gotest.tools/assert"
)

func TestBrowserParse(t *testing.T) {
	ps := NewBrowser(filepath.Join(dir, FixtureFileBrowser))
	var list []*ClientFixture
	err := ReadYamlFile(`fixtures/browser.yml`, &list)
	if err != nil {
		t.Error(err)
	}

	for _, item := range list {
		ua := item.UserAgent
		r := ps.Parse(ua)
		assert.DeepEqual(t, item.ClientMatchResult, r)
	}
}