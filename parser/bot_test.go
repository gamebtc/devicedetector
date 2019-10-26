package parser

import (
	"gotest.tools/assert"
	"path/filepath"
	"testing"
)

var botParser = NewBot(filepath.Join(dir, FixtureFileBot))

func TestGetInfoFromUABot(t *testing.T) {
	ua := `Googlebot/2.1 (http://www.googlebot.com/bot.html)`
	botParser.DiscardDetails(false)
	info := botParser.Parse(ua)
	expected := &BotMatchResult{
		Name:     `Googlebot`,
		Category: `Search bot`,
		Url:      `http://www.google.com/bot.html`,
		Producer: Producer{
			Name: "Google Inc.",
			Url:  "http://www.google.com",
		},
	}
	assert.DeepEqual(t, expected, info)
}

func TestParseNoDetails(t *testing.T) {
	ua := `Googlebot/2.1 (http://www.googlebot.com/bot.html)`
	botParser.DiscardDetails(true)
	info := botParser.Parse(ua)
	assert.Check(t, info != nil)
}

func TestParseNoBot(t *testing.T) {
	ua := `Mozilla/4.0 (compatible; MSIE 9.0; Windows NT 6.1; SV1; SE 2.x)`
	botParser.DiscardDetails(false)
	info := botParser.Parse(ua)
	assert.Check(t, info == nil)
}
