package parser

import (
	"fmt"
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

func (r *Regexp) compileGoRegex() error {
		// Make some adjustments for a different regex engine than upstream matomo
		rg := r.expression
		rg = strings.Replace(rg, `/`, `\/`, -1)
		rg = strings.Replace(rg, `++`, `+`, -1)
		rg = strings.Replace(rg, `\_`, `_`, -1)
		// if we find `\_` again, the original was `\\_`,
		// so restore that so the regex engine does not attempt to escape `_`
		rg = strings.Replace(rg, `\_`, `\\_`, -1)

		str := `(?i)(?:^|[^A-Z0-9-_]|[^A-Z0-9-]_|sprd-)(?:` + rg + ")"

		var err error
		r.regexpGo, err = regexp.Compile(str)
		return err
}

func (r *Regexp) compileLibRegex() error {
		// Make some adjustments for a different regex engine than upstream matomo
		rg := r.expression
		rg = strings.Replace(rg, `/`, `\/`, -1)
		rg = strings.Replace(rg, `++`, `+`, -1)
		rg = strings.Replace(rg, `\_`, `_`, -1)
		// if we find `\_` again, the original was `\\_`,
		// so restore that so the regex engine does not attempt to escape `_`
		rg = strings.Replace(rg, `\_`, `\\_`, -1)

		str := `(?:^|[^A-Z0-9-_]|[^A-Z0-9-]_|sprd-)(?:` + rg + ")"

		var err error
		r.regexpLib, err = regexp2.Compile(str, regexp2.IgnoreCase)
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
		panic(fmt.Sprintf("Could not compile regex '%s': %v", r.expression, err))
	}
}

func (r *Regexp) MatchString(ua string) bool {
	r.MustCompile()

	if r.regexpGo != nil {
		return r.regexpGo.MatchString(ua)
	}

	match, _ := r.regexpLib.MatchString(ua)
	return match
}

func (r *Regexp) FindStringSubmatch(ua string) []string {
	r.MustCompile()

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

