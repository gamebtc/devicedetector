package parser

import "path/filepath"

var botFactory = make(map[string]func(string) BotParser)

func RegBotParser(name string, f func(string) BotParser) {
	botFactory[name] = f
}

func GetBotCreater(name string) func(string) BotParser {
	f, _ := botFactory[name]
	return f
}

func NewBotParser(dir, name string) BotParser {
	if f, ok := botFactory[name]; ok {
		return f(dir)
	}
	return nil
}

const ParserNameBot = `bot`
const FixtureFileBot = `bots.yml`

func init() {
	RegBotParser(ParserNameBot,
		func(dir string) BotParser {
			return NewBot(filepath.Join(dir, FixtureFileBot))
		})
}

func NewBot(fileName string) *Bot {
	c := &Bot{}
	c.ParserName = ParserNameBot
	if err := c.Load(fileName); err != nil {
		return nil
	}
	return c
}

// Parses a user agent for bot information
type Bot struct {
	BotParserAbstract
}
