package client

const dir = "../../regexes/client"

type ClientFixture struct {
	UserAgent          string `yaml:"user_agent"`
	*ClientMatchResult `yaml:"client"`
}
