package client

import (
	"path/filepath"
)

const ParserNameFeedReader = `feed reader`
const FixtureFileFeedReader = `feed_readers.yml`

func init() {
	RegClientParser(ParserNameFeedReader,
		func(dir string) ClientParser {
			return NewFeedReader(filepath.Join(dir, FixtureFileFeedReader))
		})
}

func NewFeedReader(fileName string) *FeedReader {
	c := &FeedReader{}
	c.ParserName = ParserNameFeedReader
	if err := c.Load(fileName); err != nil {
		return nil
	}
	return c
}

// Client parser for feed reader detection
type FeedReader struct {
	ClientParserAbstract
}

//type FeedReader struct {
//	RegularName `yaml:",inline" json:",inline"`
//	Version     string `yaml:"version" json:"version"`
//	Url         string `yaml:"url" json:"url"`
//	Type        string `yaml:"type" json:"type"`
//}
