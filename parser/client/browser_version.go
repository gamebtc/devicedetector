package client

import (
	regexp "github.com/dlclark/regexp2"
)

// Client parser for browser engine version detection
type Version struct {
	Engine string `yaml:"engine" json:"engine"`
	regexp *regexp.Regexp
}

func (r *Version) Compile() {
	if r.regexp == nil {
		reg := r.Engine + `\s*\/?\s*((?(?=\d+\.\d)\d+[.\d]*|\d{1,7}(?=(?:\D|$))))`
		r.regexp = regexp.MustCompile(reg, regexp.IgnoreCase)
	}
}

func (r *Version) Parse(ua string) string {
	if r.regexp == nil {
		return ""
	}
	matches, _ := r.regexp.FindStringMatch(ua)
	if matches == nil {
		return ""
	}
	last := len(matches.Groups()) - 1
	if last < 0 {
		return ""
	}
	return matches.Groups()[last].String()
}
