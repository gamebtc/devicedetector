package devicedetector

import (
	"path/filepath"
	"strings"

	regexp "github.com/dlclark/regexp2"
	gover "github.com/mcuadros/go-version"

	. "github.com/gamebtc/devicedetector/parser"
	"github.com/gamebtc/devicedetector/parser/client"
	"github.com/gamebtc/devicedetector/parser/device"
)

const UNKNOWN = "UNK"
const VERSION = `3.12.1`

var desktopOsArray = []string{
	`AmigaOS`,
	`IBM`,
	`GNU/Linux`,
	`Mac`,
	`Unix`,
	`Windows`,
	`BeOS`,
	`Chrome OS`,
}

var(
	chrMobReg = regexp.MustCompile(fixUserAgentRegEx(`Chrome/[\.0-9]* Mobile`), regexp.IgnoreCase)
	chrTabReg = regexp.MustCompile(fixUserAgentRegEx(`Chrome/[\.0-9]* (?!Mobile)`), regexp.IgnoreCase)
	opaTabReg = regexp.MustCompile(fixUserAgentRegEx(`Opera Tablet`), regexp.IgnoreCase)
	opaTvReg = regexp.MustCompile(fixUserAgentRegEx(`Opera TV Store`), regexp.IgnoreCase)
)

func fixUserAgentRegEx(regex string) string {
	reg := strings.ReplaceAll(regex, `/`, `\/`)
	reg = strings.ReplaceAll(reg, `++`, `+`)
	return `(?:^|[^A-Z_-])(?:` + reg + `)`
}

type DeviceDetector struct {
	deviceParsers         []device.DeviceParser
	clientParsers         []client.ClientParser
	botParsers            []BotParser
	osParsers             []OsParser
	vendorParser          *VendorFragments
	DiscardBotInformation bool
	SkipBotDetection      bool
}

func NewDeviceDetector(dir string) (*DeviceDetector, error) {
	vp, err := NewVendor(filepath.Join(dir, FixtureFileVendor))
	if err != nil {
		return nil, err
	}
	osp, err := NewOss(filepath.Join(dir, FixtureFileOs))
	if err != nil {
		return nil, err
	}

	d := &DeviceDetector{
		vendorParser: vp,
		osParsers:    []OsParser{osp},
	}

	clientDir := filepath.Join(dir, "client")
	d.clientParsers = client.NewClientParsers(clientDir,
		[]string{
			client.ParserNameFeedReader,
			client.ParserNameMobileApp,
			client.ParserNameMediaPlayer,
			client.ParserNamePim,
			client.ParserNameBrowser,
			client.ParserNameLibrary,
		})

	deviceDir := filepath.Join(dir, "device")
	d.deviceParsers = device.NewDeviceParsers(deviceDir,
		[]string{
			device.ParserNameHbbTv,
			device.ParserNameConsole,
			device.ParserNameCar,
			device.ParserNameCamera,
			device.ParserNamePortableMediaPlayer,
			device.ParserNameMobile,
		})

	d.botParsers = []BotParser{
		NewBot(filepath.Join(dir, FixtureFileBot)),
	}


	return d, nil
}

func (d *DeviceDetector) AddClientParser(cp client.ClientParser) {
	d.clientParsers = append(d.clientParsers, cp)
}

func (d *DeviceDetector) GetClientParser() []client.ClientParser {
	return d.clientParsers
}

func (d *DeviceDetector) AddDeviceParser(dp device.DeviceParser) {
	d.deviceParsers = append(d.deviceParsers, dp)
}

func (d *DeviceDetector) GetDeviceParsers() []device.DeviceParser {
	return d.deviceParsers
}

func (d *DeviceDetector) AddBotParser(op BotParser) {
	d.botParsers = append(d.botParsers, op)
}

func (d *DeviceDetector) GetBotParsers() []BotParser {
	return d.botParsers
}

func (d *DeviceDetector) ParseBot(ua string) *BotMatchResult {
	if !d.SkipBotDetection {
		for _, parser := range d.botParsers {
			parser.DiscardDetails(d.DiscardBotInformation)
			if r := parser.Parse(ua); r != nil {
				return r
			}
		}
	}
	return nil
}

func (d *DeviceDetector) ParseOs(ua string) *OsMatchResult {
	for _, p := range d.osParsers {
		if r := p.Parse(ua); r != nil {
			return r
		}
	}
	return nil
}

func (d *DeviceDetector) ParseClient(ua string) *client.ClientMatchResult {
	for _, parser := range d.clientParsers {
		if r := parser.Parse(ua); r != nil {
			return r
		}
	}
	return nil
}

func (d *DeviceDetector) ParseDevice(ua string) *device.DeviceMatchResult {
	for _, parser := range d.deviceParsers {
		if r := parser.Parse(ua); r != nil {
			return r
		}
	}
	return nil
}

func (d *DeviceDetector) parseInfo(info *DeviceInfo) {
	ua := info.userAgent
	if r := d.ParseDevice(ua); r != nil {
		info.Type = r.Type
		info.Model = r.Model
		info.Brand = r.Brand
	}
	// If no brand has been assigned try to match by known vendor fragments
	if info.Brand == "" && d.vendorParser != nil {
		info.Brand = d.vendorParser.Parse(ua)
	}

	os := info.GetOs()
	osShortName := os.ShortName
	osFamily := GetOsFamily(osShortName)
	osVersion := os.Version
	clientName := info.GetClient().Name

	if info.Brand == "" && (osShortName == `ATV` || osShortName == `IOS` || osShortName == `MAC`) {
		info.Brand = `AP`
	}

	deviceType := GetDeviceType(info.Type)
	// Chrome on Android passes the device type based on the keyword 'Mobile'
	// If it is present the device should be a smartphone, otherwise it's a tablet
	// See https://developer.chrome.com/multidevice/user-agent#chrome_for_android_user_agent
	if deviceType == DEVICE_TYPE_INVALID && osFamily == `Android` && (clientName == `Chrome` || clientName == `Chrome Mobile`) {
		if ok, _ := chrMobReg.MatchString(ua); ok {
			deviceType = DEVICE_TYPE_SMARTPHONE
		} else if ok, _ = chrTabReg.MatchString(ua); ok {
			deviceType = DEVICE_TYPE_TABLET
		}
	}

	if deviceType == DEVICE_TYPE_INVALID {
		if info.HasAndroidMobileFragment() {
			deviceType = DEVICE_TYPE_TABLET
		} else if ok, _ := opaTabReg.MatchString(ua); ok {
			deviceType = DEVICE_TYPE_TABLET
		} else if info.HasAndroidMobileFragment() {
			deviceType = DEVICE_TYPE_SMARTPHONE
		} else if osShortName == "AND" && osVersion != "" {
			if gover.CompareSimple(osVersion, `2.0`) == -1 {
				deviceType = DEVICE_TYPE_SMARTPHONE
			} else if gover.CompareSimple(osVersion, `3.0`) >= 0 &&
				gover.CompareSimple(osVersion, `4.0`) == -1 {
				deviceType = DEVICE_TYPE_TABLET
			}
		}
	}

	// All detected feature phones running android are more likely a smartphone
	if deviceType == DEVICE_TYPE_FEATURE_PHONE && osFamily == `Android` {
		deviceType = DEVICE_TYPE_SMARTPHONE
	}

	// According to http://msdn.microsoft.com/en-us/library/ie/hh920767(v=vs.85).aspx
	if deviceType == DEVICE_TYPE_INVALID &&
		(osShortName == `WRT` || (osShortName == `WIN` && gover.CompareSimple(osVersion, `8`) >= 0)) &&
		info.IsTouchEnabled() {
		deviceType = DEVICE_TYPE_TABLET
	}

	// All devices running Opera TV Store are assumed to be a tv
	if ok, _ := opaTvReg.MatchString(ua); ok {
		deviceType = DEVICE_TYPE_TV
	}

	// Devices running Kylo or Espital TV Browsers are assumed to be a TV
	if deviceType == DEVICE_TYPE_INVALID {
		if clientName == `Kylo` || clientName == `Espial TV Browser` {
			deviceType = DEVICE_TYPE_TV
		} else if info.IsDesktop() {
			deviceType = DEVICE_TYPE_DESKTOP
		}
	}

	if deviceType != DEVICE_TYPE_INVALID {
		info.Type = GetDeviceName(deviceType)
	}
	return
}

func (d *DeviceDetector) Parse(ua string) *DeviceInfo {
	// skip parsing for empty useragents or those not containing any letter
	if !StringContainsLetter(ua) {
		return nil
	}

	info := &DeviceInfo{
		userAgent: ua,
	}

	info.bot = d.ParseBot(ua)
	if info.IsBot() {
		return info
	}

	info.os = d.ParseOs(ua)

	// Parse Clients
	// Clients might be browsers, Feed Readers, Mobile Apps, Media Players or
	// any other application accessing with an parseable UA
	info.client = d.ParseClient(ua)

	d.parseInfo(info)

	return info
}

