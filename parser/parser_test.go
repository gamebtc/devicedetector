package parser

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
)

const dir = "../regexes"
func TestVendors(t *testing.T) {
	if v, err := NewVendor(filepath.Join(dir, FixtureFileVendor)); err != nil {
		t.Fatal(err)
	} else if v == nil {
		t.Fatal("value is null")
	} else {
		str, _ := json.Marshal(v)
		fmt.Printf(string(str))
	}
}

func TestReg(t *testing.T) {
	name := `Chrome(?:/(\d+[\.\d]+))?`
	ua := `Mozilla/5.0 (Linux; Android 4.2.2; ARCHOS 101 PLATINUM Build/JDQ39) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.114 Safari/537.36`
	matches := MatchUserAgent(ua, name)
	for _, m := range matches {
		fmt.Print(m)
		fmt.Print("\r\n")
	}
}
