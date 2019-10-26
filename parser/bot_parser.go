package parser

import "strings"

type Producer struct {
	Name string `yaml:"name" json:"name"`
	Url  string `yaml:"url" json:"url"`
}

type BotMatchResult struct {
	Name     string   `yaml:"name" json:"name"`
	Category string   `yaml:"category" json:"category"`
	Url      string   `yaml:"url" json:"url"`
	Producer Producer `yaml:"producer" json:"producer"`
}

func (b *BotMatchResult) Equal(a *BotMatchResult) bool {
	return b.Name == a.Name &&
		b.Category == a.Category &&
		b.Url == a.Url &&
		b.Producer.Name == a.Producer.Name &&
		b.Producer.Url == a.Producer.Url
}

type BotReg struct {
	Regular        `yaml:",inline" json:",inline"`
	BotMatchResult `yaml:",inline" json:",inline"`
}

type BotParser interface {
	PreMatch(string) bool
	Parse(string) *BotMatchResult
	DiscardDetails(bool)
}

// Abstract class for all bot parsers
type BotParserAbstract struct {
	Regexes        []*BotReg
	ParserName     string
	discardDetails bool
	overAllMatch   Regular
}

func (b *BotParserAbstract) DiscardDetails(v bool) {
	b.discardDetails = v
}

func (b *BotParserAbstract) Load(file string) error {
	var v []*BotReg
	err := ReadYamlFile(file, &v)
	if err != nil {
		return err
	}
	for _, item := range v {
		item.Compile()
	}
	b.Regexes = v
	return nil
}

func (b *BotParserAbstract) PreMatch(ua string) bool {
	if b.overAllMatch.Regexp == nil {
		count := len(b.Regexes)
		if count == 0 {
			return false
		}
		sb := strings.Builder{}
		sb.WriteString(b.Regexes[count-1].Regex)
		for i := count - 2; i >= 0; i-- {
			sb.WriteString("|")
			sb.WriteString(b.Regexes[i].Regex)
		}
		b.overAllMatch.Regex = sb.String()
		b.overAllMatch.Compile()
	}
	r := b.overAllMatch.IsMatchUserAgent(ua)
	return r
}

var EmptyBotMatchResult = new(BotMatchResult)

// Parses the current UA and checks whether it contains bot information
func (b *BotParserAbstract) Parse(ua string) *BotMatchResult {
	if b.PreMatch(ua) {
		if b.discardDetails {
			return EmptyBotMatchResult
		}
		for _, regex := range b.Regexes {
			matches := regex.MatchUserAgent(ua)
			if len(matches) > 0 {
				return &regex.BotMatchResult
			}
		}
	}
	return nil
}
