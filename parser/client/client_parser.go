package client

import (
	"sort"
	"strings"

	. "github.com/muxinc/devicedetector/parser"
)

type ClientMatchResult struct {
	Type    string `yaml:"type" json:"type"`
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`

	ShortName     string `yaml:"short_name" json:"short_name"`
	Engine        string `yaml:"engine" json:"engine"`
	EngineVersion string `yaml:"engine_version" json:"engine_version"`
}

type ClientParser interface {
	PreMatch(string) bool
	Parse(string) *ClientMatchResult
}

type ClientReg struct {
	Regular `yaml:",inline" json:",inline"`
	Name    string `yaml:"name" json:"name"`
	Version string `yaml:"version" json:"version"`
}

// Parses the current UA and checks whether it contains any client information
type ClientParserAbstract struct {
	Regexes      []*ClientReg
	ParserName   string
	overAllMatch Regular
}

func (c *ClientParserAbstract) Load(file string) error {
	var v []*ClientReg
	err := ReadYamlFile(file, &v)
	if err != nil {
		return err
	}
	for _, item := range v {
		item.Compile()
	}
	c.Regexes = v
	return nil
}

func (c *ClientParserAbstract) PreMatch(ua string) bool {
	if c.overAllMatch.Regexp == nil {
		count := len(c.Regexes)
		if count == 0 {
			return false
		}
		sb := strings.Builder{}
		sb.WriteString(c.Regexes[count-1].Regex)
		for i := count - 2; i >= 0; i-- {
			sb.WriteString("|")
			sb.WriteString(c.Regexes[i].Regex)
		}
		c.overAllMatch.Regex = sb.String()
		c.overAllMatch.Compile()
	}
	r := c.overAllMatch.IsMatchUserAgent(ua)
	return r
}

// Parses the current UA and checks whether it contains any client information
func (c *ClientParserAbstract) Parse(ua string) *ClientMatchResult {
	if c.PreMatch(ua) {
		for _, regex := range c.Regexes {
			matches := regex.MatchUserAgent(ua)
			if len(matches) > 0 {
				return &ClientMatchResult{
					Type:    c.ParserName,
					Name:    BuildByMatch(regex.Name, matches),
					Version: BuildVersion(regex.Version, matches),
				}
			}
		}
	}
	return nil
}

// Returns all names defined in the regexes
// Attention: This method might not return all names of detected clients
func (c *ClientParserAbstract) GetAvailableClients() []string {
	names := make([]string, 0, len(c.Regexes))
	for _, regex := range c.Regexes {
		n := regex.Name
		if n != `$1` && ArrayContainsString(names, n) == false {
			names = append(names, n)
		}
	}
	sort.Strings(names)
	return names
}
