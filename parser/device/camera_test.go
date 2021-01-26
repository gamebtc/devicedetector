package device

import (
	"path/filepath"
	"testing"

	. "github.com/muxinc/devicedetector/parser"
	"gotest.tools/assert"
)

func TestCameraParse(t *testing.T) {
	ps := NewCamera(filepath.Join(dir, FixtureFileCamera))
	var list []*DeviceFixture
	err := ReadYamlFile(`fixtures/camera.yml`, &list)
	if err != nil {
		t.Error(err)
	}

	//ua := `Mozilla/5.0 (Linux; U; Android 4.0; de-DE; EK-GC100 Build/IMM76D) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30`
	//item := DeviceFixture{
	//	UserAgent: ua,
	//	DeviceFixtureResult: DeviceFixtureResult{
	//		Type:  8,
	//		Brand: `SA`,
	//		Model: `GALAXY Camera`,
	//	},
	//}
	//r := ps.Parse(ua)
	//test := item.GetDeviceMatchResult()
	//assert.DeepEqual(t, test, r)

	for _, item := range list {
		ua := item.UserAgent
		r := ps.Parse(ua)
		test := item.GetDeviceMatchResult()
		assert.DeepEqual(t, test, r)
	}
}
