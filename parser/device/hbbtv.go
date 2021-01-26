package device

import (
	"path/filepath"

	. "github.com/muxinc/devicedetector/parser"
)

const ParserNameHbbTv = `tv`
const FixtureFileHbbTv = `televisions.yml`

func init() {
	RegDeviceParser(ParserNameHbbTv,
		func(dir string) DeviceParser {
			return NewHbbTv(filepath.Join(dir, FixtureFileHbbTv))
		})
}

func NewHbbTv(fileName string) *HbbTv {
	h := &HbbTv{}
	if err := h.Load(fileName); err != nil {
		h.hbbTvRegx.Regex = `HbbTV/([1-9]{1}(?:.[0-9]{1}){1,2})`
		return nil
	}
	return h
}

// Device parser for hbbtv detection
type HbbTv struct {
	DeviceParserAbstract
	hbbTvRegx Regular
}

func (h *HbbTv) Parse(ua string) *DeviceMatchResult {
	// only parse user agents containing hbbtv fragment
	if !h.IsHbbTv(ua) {
		return nil
	}
	r := h.DeviceParserAbstract.Parse(ua)
	// always set device type to tv, even if no model/brand could be found
	if r != nil {
		r.Type = ParserNameHbbTv
	}
	return r
}

// Returns if the parsed UA was identified as a HbbTV device
func (h *HbbTv) IsHbbTv(ua string) bool {
	return h.hbbTvRegx.IsMatchUserAgent(ua)
}
