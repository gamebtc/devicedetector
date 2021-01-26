package client

import (
	"path/filepath"
	"testing"

	. "github.com/muxinc/devicedetector/parser"
	"gotest.tools/assert"
)

func TestFeedReaderParse(t *testing.T) {
	ps := NewFeedReader(filepath.Join(dir, FixtureFileFeedReader))
	var list []*ClientFixture
	err := ReadYamlFile(`fixtures/feed_reader.yml`, &list)
	if err != nil {
		t.Error(err)
	}

	for _, item := range list {
		ua := item.UserAgent
		r := ps.Parse(ua)
		assert.DeepEqual(t, item.ClientMatchResult, r)
	}
}
