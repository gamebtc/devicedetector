package device

import (
	"path/filepath"
	"testing"

	"gotest.tools/assert"
)

func TestHbbTvParse(t *testing.T) {
	ps := NewHbbTv(filepath.Join(dir, FixtureFileHbbTv))
	ua := `Opera/9.80 (Linux mips ; U; HbbTV/1.1.1 (; Philips; ; ; ; ) CE-HTML/1.0 NETTV/3.2.1; en) Presto/2.6.33 Version/10.70`

	r := ps.Parse(ua)
	assert.Equal(t, r.Type, ParserNameHbbTv)
}
