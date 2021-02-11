package parser

import (
	"log"
	"regexp"
	goregexp "regexp"
	"strings"

	regexp2 "github.com/dlclark/regexp2"
)

type Regexp struct {
	expression string
	regexpGo *goregexp.Regexp
	regexpLib *regexp2.Regexp
}

func NewRegexp(re string) *Regexp {
	r := &Regexp{expression: re}
	r.MustCompile()
	return r
}

// Make some adjustments for a different regex engine than upstream matomo
func cleanRegexString(re string) string {
	rg := strings.Replace(re, `/`, `\/`, -1)
	rg = strings.Replace(rg, `++`, `+`, -1)
	rg = strings.Replace(rg, `\_`, `_`, -1)
	// if we find `\_` again, the original was `\\_`,
	// so restore that so the regex engine does not attempt to escape `_`
	rg = strings.Replace(rg, `\_`, `\\_`, -1)

	// only match if useragent begins with given regex or there is no letter before it
	return `(?:^|[^A-Z0-9-_]|[^A-Z0-9-]_|sprd-)(?:` + rg + ")"
}

func (r *Regexp) compileGoRegex() error {
		var err error
		r.regexpGo, err = regexp.Compile(`(?i)` + cleanRegexString(r.expression))
		return err
}

func (r *Regexp) compileLibRegex() error {
		var err error
		r.regexpLib, err = regexp2.Compile(cleanRegexString(r.expression), regexp2.IgnoreCase)
		return err
}

func (r *Regexp) Compile() error {
	if r.regexpGo == nil && r.regexpLib == nil {
		if err := r.compileGoRegex(); err != nil {
			return r.compileLibRegex()
		}
	}
	return nil
}

func (r *Regexp) MustCompile() {
	if err := r.Compile(); err != nil {
		log.Panicf("Could not compile regex '%s': %v", r.expression, err)
	}
}

func (r *Regexp) MatchString(ua string) bool {
	if r.regexpGo != nil {
		return r.regexpGo.MatchString(ua)
	}

	match, _ := r.regexpLib.MatchString(ua)
	return match
}

func (r *Regexp) FindStringSubmatch(ua string) []string {
	if r.regexpGo != nil {
		return r.regexpGo.FindStringSubmatch(ua)
	}

	if match, err := r.regexpLib.FindStringMatch(ua); err == nil && match != nil {
		matches := make([]string, match.GroupCount())
		for i, g := range match.Groups() {
			matches[i] = g.String()
		}
		return matches
	}
	return nil
}

