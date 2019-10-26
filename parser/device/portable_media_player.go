package device

import (
	"path/filepath"
)

const ParserNamePortableMediaPlayer = `portablemediaplayer`
const FixtureFilePortableMediaPlayer = `portable_media_player.yml`

func init() {
	RegDeviceParser(ParserNamePortableMediaPlayer,
		func(dir string) DeviceParser {
			return NewPortableMediaPlayer(filepath.Join(dir, FixtureFilePortableMediaPlayer))
		})
}

func NewPortableMediaPlayer(fileName string) *PortableMediaPlayer {
	p := &PortableMediaPlayer{}
	if err := p.Load(fileName); err != nil {
		return nil
	}
	return p
}

// Device parser for portable media player detection
type PortableMediaPlayer struct {
	DeviceParserAbstract
}

func (p *PortableMediaPlayer) Parse(ua string) *DeviceMatchResult {
	if !p.PreMatch(ua) {
		return nil
	}
	return p.DeviceParserAbstract.Parse(ua)
}
