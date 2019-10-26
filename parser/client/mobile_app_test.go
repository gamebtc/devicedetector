package client

import (
	"path/filepath"
	"testing"

	. "github.com/gamebtc/devicedetector/parser"
	"gotest.tools/assert"
)

func TestMediaAppParse(t *testing.T) {
	ps := NewMobileApp(filepath.Join(dir, FixtureFileMobileApp))
	var list []*ClientFixture
	err := ReadYamlFile(`fixtures/mobile_app.yml`, &list)
	if err != nil {
		t.Error(err)
	}

	for _, item := range list {
		ua := item.UserAgent
		r := ps.Parse(ua)
		assert.DeepEqual(t, item.ClientMatchResult, r)
	}
}
