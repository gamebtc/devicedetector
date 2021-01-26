package device

import (
	"sort"
	"strings"

	. "github.com/muxinc/devicedetector/parser"
)

const UnknownBrand = "Unknown"

type DeviceMatchResult struct {
	Type  string `yaml:"type"`
	Model string `yaml:"model"`
	Brand string `yaml:"brand"`
}

type DeviceParser interface {
	PreMatch(string) bool
	Parse(string) *DeviceMatchResult
}

type Model struct {
	Regular `yaml:",inline" json:",inline"`
	Model   string `yaml:"model" json:"model"`
	Device  string `yaml:"device" json:"device"` //mobile
	Brand   string `yaml:"brand" json:"brand"`   //mobile
}

type DeviceReg struct {
	Regular `yaml:",inline" json:",inline"`
	Model   string   `yaml:"model" json:"model"`
	Device  string   `yaml:"device" json:"device"`
	Models  []*Model `yaml:"models" json:"models"`
}

type DeviceParserAbstract struct {
	Regexes      map[string]*DeviceReg
	overAllMatch Regular
}

func (d *DeviceParserAbstract) Load(file string) error {
	var v map[string]*DeviceReg
	err := ReadYamlFile(file, &v)
	if err != nil {
		return err
	}
	for _, item := range v {
		item.Compile()
		for _, m := range item.Models {
			m.Compile()
		}
	}
	d.Regexes = v
	return nil
}

func (d *DeviceParserAbstract) PreMatch(ua string) bool {
	if d.overAllMatch.Regexp == nil {
		count := len(d.Regexes)
		if count == 0 {
			return false
		}
		sortKeys := make([]string, 0, count)
		for k, _ := range d.Regexes {
			sortKeys = append(sortKeys, k)
		}
		sort.Strings(sortKeys)
		sb := strings.Builder{}
		sb.WriteString(d.Regexes[sortKeys[count-1]].Regex)
		for i := count - 2; i >= 0; i-- {
			sb.WriteString("|")
			sb.WriteString(d.Regexes[sortKeys[i]].Regex)
		}
		d.overAllMatch.Regex = sb.String()
		d.overAllMatch.Compile()
	}
	r := d.overAllMatch.IsMatchUserAgent(ua)
	return r
}

func (d *DeviceParserAbstract) Parse(ua string) *DeviceMatchResult {
	var regex *DeviceReg
	var brand string
	var matches []string
	for brand, regex = range d.Regexes {
		matches = regex.MatchUserAgent(ua)
		if len(matches) > 0 {
			break
		}
	}

	if regex == nil || len(matches) == 0 {
		return nil
	}

	r := &DeviceMatchResult{
		Type: regex.Device,
	}
	if brand != UnknownBrand {
		brandId := FindBrand(brand)
		if brandId == "" {
			return nil
		}
		r.Brand = brandId
	}

	if regex.Model != "" {
		r.Model = BuildModel(regex.Model, matches)
	}

	for _, modelRegex := range regex.Models {
		modelMatches := modelRegex.MatchUserAgent(ua)
		if len(modelMatches) > 0 {
			r.Model = strings.TrimSpace(BuildModel(modelRegex.Model, modelMatches))
			if modelRegex.Brand != "" {
				if brandId := FindBrand(modelRegex.Brand); brandId != "" {
					r.Brand = brandId
				}
			}
			if modelRegex.Device != "" {
				r.Type = modelRegex.Device
			}
			return r
		}
	}
	return r
}
