package parser

import (
	"path/filepath"
	"testing"

	"gotest.tools/assert"
)



func TestOsParse(t *testing.T) {
	type OsFixture struct {
		OsMatchResult `yaml:"os" json:"os"`
		UserAgent     string `yaml:"user_agent" json:"user_agent"`
	}

	var osParser, _ = NewOss(filepath.Join(dir, FixtureFileOs))

	var list []OsFixture
	err := ReadYamlFile(`fixtures/oss.yml`, &list)
	if err != nil {
		t.Error(err)
	}

	for _, item := range list {
		r := osParser.Parse(item.UserAgent)
		//assert.Check(t, NameEqual(item.Name, r.Name), "%v:%v_%v", item.UserAgent, item.Name, r.Name)
		//assert.Check(t, NameEqual(item.ShortName, r.ShortName), )
		//assert.Check(t, NameEqual(item.Version, r.Version), )
		//assert.Check(t, NameEqual(item.Platform, r.Platform), )
		assert.DeepEqual(t, &item.OsMatchResult, r)
	}
}

func TestOsInGroup(t *testing.T) {
	for os, _ := range OperatingSystems {
		contains := false
		for _, familyOs := range OsFamilies {
			if ArrayContainsString(familyOs, os) {
				contains = true
				break
			}
		}
		assert.Check(t, contains)
	}
}

func TestGetNameFromId(t *testing.T) {
	data := [][]string{
		[]string{"DEB", "4.5", "Debian 4.5"},
		[]string{"WRT", "", "Windows RT"},
		[]string{"WIN", "98", "Windows 98"},
		[]string{"XXX", "4.5", ""},
	}
	for _, item := range data {
		os, version, expected := item[0], item[1], item[2]
		r := GetOsNameFromId(os, version)
		assert.Equal(t, r, expected)
	}
}
