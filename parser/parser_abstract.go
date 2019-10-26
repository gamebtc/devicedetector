package parser

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"

	regexp "github.com/dlclark/regexp2"
)

const Unknown = "Unknown"
const UnknownShort = "UNK"

// Indicates how deep versioning will be detected
// if $maxMinorParts is 0 only the major version will be returned
var maxMinorParts = -1

// Versioning constant used to set max versioning to major version only
// Version examples are: 3, 5, 6, 200, 123, ...
const VERSION_TRUNCATION_MAJOR = 0

// Versioning constant used to set max versioning to minor version
// Version examples are: 3.4, 5.6, 6.234, 0.200, 1.23, ...
const VERSION_TRUNCATION_MINOR = 1

// Versioning constant used to set max versioning to path level
// Version examples are: 3.4.0, 5.6.344, 6.234.2, 0.200.3, 1.2.3, ...
const VERSION_TRUNCATION_PATCH = 2

// Versioning constant used to set versioning to build number
// Version examples are: 3.4.0.12, 5.6.334.0, 6.234.2.3, 0.200.3.1, 1.2.3.0, ...
const VERSION_TRUNCATION_BUILD = 3

// Versioning constant used to set versioning to unlimited (no truncation)
const VERSION_TRUNCATION_NONE = -1

func SetVersionTruncation(t int) {
	if t == VERSION_TRUNCATION_BUILD ||
		t == VERSION_TRUNCATION_NONE ||
		t == VERSION_TRUNCATION_MAJOR ||
		t == VERSION_TRUNCATION_MINOR ||
		t == VERSION_TRUNCATION_PATCH {
		maxMinorParts = t
	}
}

func ArrayContainsString(list []string, v string) bool {
	if list != nil {
		for _, i := range list {
			if i == v {
				return true
			}
		}
	}
	return false
}

func ReadYamlFile(file string, v interface{}) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.New("not exists:" + file)
	}
	return yaml.Unmarshal(data, v)
}

type MatchResult interface {
	GetName() string
	SetName(string)
}

type Regular struct {
	Regex  string `yaml:"regex" json:"regex"`
	Regexp *regexp.Regexp
}

func (r *Regular) Compile() *regexp.Regexp {
	if r.Regexp == nil {
		// $regex = '/(?:^|[^A-Z_-])(?:' . str_replace('/', '\/', $regex) . ')/i';
		//str := `(?i)(?:^|[^A-Z0-9-_]|[^A-Z0-9-]_|sprd-)(?:` + r.Regex + ")"
		rg := r.Regex
		rg = strings.Replace(rg, `/`, `\/`, -1)
		rg = strings.Replace(rg, `++`, `+`, -1)
		rg = strings.Replace(rg, `\_`, `_`, -1)
		str := `(?:^|[^A-Z0-9-_]|[^A-Z0-9-]_|sprd-)(?:` + rg + ")"
		r.Regexp = regexp.MustCompile(str, regexp.IgnoreCase)
	}
	return r.Regexp
}

func (r *Regular) IsMatchUserAgent(ua string) bool {
	m, _ := r.Compile().MatchString(ua)
	return m
}

func (r *Regular) MatchUserAgent(ua string) []string {
	//return r.Compile().FindStringSubmatch(ua)
	if match, err := r.Compile().FindStringMatch(ua); err == nil && match != nil {
		matches := make([]string, match.GroupCount())
		for i, g := range match.Groups() {
			matches[i] = g.String()
		}
		return matches
	}
	return nil
}

func MatchUserAgent(ua, regex string) []string {
	rx := Regular{Regex: regex}
	return rx.MatchUserAgent(ua)
}

func BuildByMatch(item string, matches []string) string {
	if strings.IndexByte(item, byte('$')) == -1 {
		return item
	}
	for nb := 1; nb <= 3; nb++ {
		key := "$" + strconv.Itoa(nb)
		if !strings.Contains(item, key) {
			continue
		}
		replace := ""
		if nb < len(matches) {
			replace = matches[nb]
		}
		item = strings.TrimSpace(strings.ReplaceAll(item, key, replace))
	}
	return item
}

func BuildVersion(versionString string, matches []string) string {
	ver := BuildByMatch(versionString, matches)
	ver = strings.TrimRight(strings.ReplaceAll(ver, "_", "."), ".")
	verParts := strings.Split(ver, ".")
	if maxMinorParts == -1 || len(verParts)-1 <= maxMinorParts {
		return ver
	}
	newVerParts := make([]string, 1+maxMinorParts)
	copy(newVerParts, verParts)
	ver = strings.Join(newVerParts, ".")
	return ver
}

var(
	tdReg = regexp.MustCompile(` TD$`, regexp.IgnoreCase)
)
func BuildModel(m string, matches []string) string {
	model := BuildByMatch(m, matches)
	model = strings.ReplaceAll(model, "_", " ")
	model, _ = tdReg.Replace(model, "", 0, -1)
	if model == "Build" {
		return ""
	}
	return model
}

func StringEqualIgnoreCase(a string, b string) bool {
	l := len(a)
	if l != len(b) {
		return false
	}
	const toLowerChar = 'a' - 'A'
	for i := 0; i < l; i++ {
		ac := a[i]
		bc := b[i]
		if ac != bc {
			if ac >= 'A' && ac <= 'Z' {
				ac += toLowerChar
			}
			if bc >= 'A' && bc <= 'Z' {
				bc += toLowerChar
			}
			if ac != bc {
				return false
			}
		}
	}
	return true
}

func NameEqual(a string, b string) bool {
	if StringEqualIgnoreCase(a, b) {
		return true
	}
	if a == UnknownShort || a == Unknown {
		a = ""
	}
	if b == UnknownShort || b == Unknown {
		b = ""
	}
	return a == b
}

func StringContainsLetter(ua string) bool {
	for _, c := range ua {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			return true
		}
	}
	return false
}