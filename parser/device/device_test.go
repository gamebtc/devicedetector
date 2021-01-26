package device

import(
	. "github.com/muxinc/devicedetector/parser"
)

const dir = "../../regexes/device"

type DeviceFixtureResult struct {
	Type  int `yaml:"type"`
	Model string `yaml:"model"`
	Brand string `yaml:"brand"`
}

type DeviceFixture struct {
	UserAgent           string `yaml:"user_agent"`
	DeviceFixtureResult `yaml:"device"`
}

func (d *DeviceFixture)GetDeviceMatchResult()*DeviceMatchResult {
	return &DeviceMatchResult{
		Model: d.Model,
		Brand: d.Brand,
		Type:  GetDeviceName(d.Type),
	}
}
