package devicedetector

import (
	regexp "github.com/dlclark/regexp2"

	. "github.com/gamebtc/devicedetector/parser"
	"github.com/gamebtc/devicedetector/parser/client"
	"github.com/gamebtc/devicedetector/parser/device"
)

var (
	touchReg  = regexp.MustCompile(fixUserAgentRegEx(`Touch`), regexp.IgnoreCase)
	adrTabReg = regexp.MustCompile(fixUserAgentRegEx(`Android( [\.0-9]+)?; Tablet;`), regexp.IgnoreCase)
	adrMobReg = regexp.MustCompile(fixUserAgentRegEx(`Android( [\.0-9]+)?; Mobile;`), regexp.IgnoreCase)
)

type DeviceInfo struct {
	userAgent string
	device.DeviceMatchResult
	client *client.ClientMatchResult
	os     *OsMatchResult
	bot    *BotMatchResult
}

func (d *DeviceInfo) GetDeviceType() int {
	return GetDeviceType(d.Type)
}

func (d *DeviceInfo) IsBot() bool {
	return d.bot != nil
}

func (d *DeviceInfo) IsTouchEnabled() bool {
	find, _ := touchReg.MatchString(d.userAgent)
	return find
}

func (d *DeviceInfo) HasAndroidTableFragment() bool {
	find, _ := adrTabReg.MatchString(d.userAgent)
	return find
}

func (d *DeviceInfo) HasAndroidMobileFragment() bool {
	find, _ := adrMobReg.MatchString(d.userAgent)
	return find
}

func (d *DeviceInfo) UsesMobileBrowser() bool {
	return d.client != nil && d.client.Type == client.ParserNameBrowser && client.IsMobileOnlyBrowser(d.client.ShortName)
}

func (d *DeviceInfo) IsMobile() bool {
	if d.Type != "" {
		if deviceType := GetDeviceType(d.Type); DEVICE_TYPE_INVALID != deviceType {
			switch deviceType {
			case DEVICE_TYPE_FEATURE_PHONE,
				DEVICE_TYPE_SMARTPHONE,
				DEVICE_TYPE_TABLET,
				DEVICE_TYPE_PHABLET,
				DEVICE_TYPE_CAMERA,
				DEVICE_TYPE_PORTABLE_MEDIA_PAYER:
				return true
			case DEVICE_TYPE_TV,
				DEVICE_TYPE_SMART_DISPLAY,
				DEVICE_TYPE_CONSOLE:
				return false
			}
		}
	}

	if d.UsesMobileBrowser() {
		return true
	}

	if d.os == nil || d.os.ShortName == "" || d.os.ShortName == UNKNOWN {
		return false
	}

	return !d.IsBot() && !d.IsDesktop()
}

func (d *DeviceInfo) IsDesktop() bool {
	if d.os == nil || d.os.ShortName == "" || d.os.ShortName == UNKNOWN {
		return false
	}

	// Check for browsers available for mobile devices only
	if d.UsesMobileBrowser() {
		return false
	}

	if decodedFamily := GetOsFamily(d.os.ShortName); decodedFamily != "" {
		return ArrayContainsString(desktopOsArray, decodedFamily)
	}
	return false
}

func (d *DeviceInfo) GetOs() *OsMatchResult {
	if d.os != nil {
		return d.os
	}
	return &OsMatchResult{}
}

func (d *DeviceInfo) GetClient() *client.ClientMatchResult {
	if d.client != nil {
		return d.client
	}
	return &client.ClientMatchResult{}
}

func (d *DeviceInfo) GetBrowserClient() *client.BrowserMatchResult {
	if d.client != nil && d.client.Type == client.ParserNameBrowser {
		return d.client
	}
	return &client.BrowserMatchResult{}
}

func (d *DeviceInfo) GetDevice() *device.DeviceMatchResult {
	return &d.DeviceMatchResult
}

func (d *DeviceInfo) GetDeviceName() string {
	return d.Type
}

func (d *DeviceInfo) GetBrand() string {
	return d.Brand
}

func (d *DeviceInfo) GetBrandName() string {
	return GetFullName(d.Brand)
}

func (d *DeviceInfo) GetModel() string {
	return d.Model
}

func (d *DeviceInfo) GetUserAgent() string {
	return d.userAgent
}

func (d *DeviceInfo) GetBot() *BotMatchResult {
	return d.bot
}

func (d *DeviceInfo) GetOsFamily() string {
	if d.os != nil {
		return GetOsFamily(d.os.ShortName)
	}
	return ""
}

func (d *DeviceInfo) GetBrowserFamily() string {
	if d.client != nil {
		if v, ok := client.GetBrowserFamily(d.client.ShortName); ok {
			return v
		}
	}
	return ""
}
