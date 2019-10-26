package client

import (
	"path/filepath"
	
	gover "github.com/mcuadros/go-version"

	. "github.com/gamebtc/devicedetector/parser"
)

// Known browsers mapped to their internal short codes
var availableBrowsers = map[string]string{
	`2B` : `2345 Browser`,
	`36` : `360 Phone Browser`,
	`3B` : `360 Browser`,
	`AA` : `Avant Browser`,
	`AB` : `ABrowse`,
	`AF` : `ANT Fresco`,
	`AG` : `ANTGalio`,
	`AL` : `Aloha Browser`,
	`AM` : `Amaya`,
	`AO` : `Amigo`,
	`AN` : `Android Browser`,
	`AD` : `AOL Shield`,
	`AR` : `Arora`,
	`AV` : `Amiga Voyager`,
	`AW` : `Amiga Aweb`,
	`AT` : `Atomic Web Browser`,
	`AS` : `Avast Secure Browser`,
	`BA` : `Beaker Browser`,
	`BB` : `BlackBerry Browser`,
	`BD` : `Baidu Browser`,
	`BS` : `Baidu Spark`,
	`BI` : `Basilisk`,
	`BE` : `Beonex`,
	`BJ` : `Bunjalloo`,
	`BL` : `B-Line`,
	`BR` : `Brave`,
	`BK` : `BriskBard`,
	`BX` : `BrowseX`,
	`CA` : `Camino`,
	`CC` : `Coc Coc`,
	`CD` : `Comodo Dragon`,
	`C1` : `Coast`,
	`CX` : `Charon`,
	`CE` : `CM Browser`,
	`CF` : `Chrome Frame`,
	`HC` : `Headless Chrome`,
	`CH` : `Chrome`,
	`CI` : `Chrome Mobile iOS`,
	`CK` : `Conkeror`,
	`CM` : `Chrome Mobile`,
	`CN` : `CoolNovo`,
	`CO` : `CometBird`,
	`CP` : `ChromePlus`,
	`CR` : `Chromium`,
	`CY` : `Cyberfox`,
	`CS` : `Cheshire`,
	`CU` : `Cunaguaro`,
	`CV` : `Chrome Webview`,
	`DB` : `dbrowser`,
	`DE` : `Deepnet Explorer`,
	`DF` : `Dolphin`,
	`DO` : `Dorado`,
	`DL` : `Dooble`,
	`DI` : `Dillo`,
	`DD` : `DuckDuckGo Privacy Browser`,
	`EC` : `Ecosia`,
	`EI` : `Epic`,
	`EL` : `Elinks`,
	`EB` : `Element Browser`,
	`EP` : `GNOME Web`,
	`ES` : `Espial TV Browser`,
	`F1` : `Firefox Mobile iOS`,
	`FB` : `Firebird`,
	`FD` : `Fluid`,
	`FE` : `Fennec`,
	`FF` : `Firefox`,
	`FK` : `Firefox Focus`,
	`FR` : `Firefox Rocket`,
	`FL` : `Flock`,
	`FM` : `Firefox Mobile`,
	`FW` : `Fireweb`,
	`FN` : `Fireweb Navigator`,
	`FU` : `FreeU`,
	`GA` : `Galeon`,
	`GE` : `Google Earth`,
	`HA` : `Hawk Turbo Browser`,
	`HO` : `hola! Browser`,
	`HJ` : `HotJava`,
	`HU` : `Huawei Browser`,
	`IB` : `IBrowse`,
	`IC` : `iCab`,
	`I2` : `iCab Mobile`,
	`I1` : `Iridium`,
	`I3` : `Iron Mobile`,
	`I4` : `IceCat`,
	`ID` : `IceDragon`,
	`IV` : `Isivioo`,
	`IW` : `Iceweasel`,
	`IE` : `Internet Explorer`,
	`IM` : `IE Mobile`,
	`IR` : `Iron`,
	`JS` : `Jasmine`,
	`JI` : `Jig Browser`,
	`JO` : `Jio Browser`,
	`KB` : `K.Browser`,
	`KI` : `Kindle Browser`,
	`KM` : `K-meleon`,
	`KO` : `Konqueror`,
	`KP` : `Kapiko`,
	`KW` : `Kiwi`,
	`KY` : `Kylo`,
	`KZ` : `Kazehakase`,
	`LB` : `Cheetah Browser`,
	`LF` : `LieBaoFast`,
	`LG` : `LG Browser`,
	`LI` : `Links`,
	`LU` : `LuaKit`,
	`LS` : `Lunascape`,
	`LX` : `Lynx`,
	`MB` : `MicroB`,
	`MC` : `NCSA Mosaic`,
	`ME` : `Mercury`,
	`MF` : `Mobile Safari`,
	`MI` : `Midori`,
	`MO` : `Mobicip`,
	`MU` : `MIUI Browser`,
	`MS` : `Mobile Silk`,
	`MT` : `Mint Browser`,
	`MX` : `Maxthon`,
	`NB` : `Nokia Browser`,
	`NO` : `Nokia OSS Browser`,
	`NV` : `Nokia Ovi Browser`,
	`NX` : `Nox Browser`,
	`NE` : `NetSurf`,
	`NF` : `NetFront`,
	`NL` : `NetFront Life`,
	`NP` : `NetPositive`,
	`NS` : `Netscape`,
	`NT` : `NTENT Browser`,
	`OC` : `Oculus Browser`,
	`O1` : `Opera Mini iOS`,
	`OB` : `Obigo`,
	`OD` : `Odyssey Web Browser`,
	`OF` : `Off By One`,
	`OE` : `ONE Browser`,
	`OG` : `Opera Neon`,
	`OH` : `Opera Devices`,
	`OI` : `Opera Mini`,
	`OM` : `Opera Mobile`,
	`OP` : `Opera`,
	`ON` : `Opera Next`,
	`OO` : `Opera Touch`,
	`OR` : `Oregano`,
	`OV` : `Openwave Mobile Browser`,
	`OW` : `OmniWeb`,
	`OT` : `Otter Browser`,
	`PL` : `Palm Blazer`,
	`PM` : `Pale Moon`,
	`PP` : `Oppo Browser`,
	`PR` : `Palm Pre`,
	`PU` : `Puffin`,
	`PW` : `Palm WebPro`,
	`PA` : `Palmscape`,
	`PX` : `Phoenix`,
	`PO` : `Polaris`,
	`PT` : `Polarity`,
	`PS` : `Microsoft Edge`,
	`Q1` : `QQ Browser Mini`,
	`QQ` : `QQ Browser`,
	`QT` : `Qutebrowser`,
	`QZ` : `QupZilla`,
	`QM` : `Qwant Mobile`,
	`QW` : `QtWebEngine`,
	`RE` : `Realme Browser`,
	`RK` : `Rekonq`,
	`RM` : `RockMelt`,
	`SB` : `Samsung Browser`,
	`SA` : `Sailfish Browser`,
	`SC` : `SEMC-Browser`,
	`SE` : `Sogou Explorer`,
	`SF` : `Safari`,
	`SH` : `Shiira`,
	`SK` : `Skyfire`,
	`SS` : `Seraphic Sraf`,
	`SL` : `Sleipnir`,
	`SN` : `Snowshoe`,
	`SO` : `Sogou Mobile Browser`,
	`SI` : `Sputnik Browser`,
	`SR` : `Sunrise`,
	`SP` : `SuperBird`,
	`ST` : `Streamy`,
	`SX` : `Swiftfox`,
	`SZ` : `Seznam Browser`,
	`TF` : `TenFourFox`,
	`TB` : `Tenta Browser`,
	`TZ` : `Tizen Browser`,
	`TS` : `TweakStyle`,
	`UC` : `UC Browser`,
	`UM` : `UC Browser Mini`,
	`VI` : `Vivaldi`,
	`VV` : `vivo Browser`,
	`VB` : `Vision Mobile Browser`,
	`WP` : `Web Explorer`,
	`WE` : `WebPositive`,
	`WF` : `Waterfox`,
	`WH` : `Whale Browser`,
	`WO` : `wOSBrowser`,
	`WT` : `WeTab Browser`,
	`YA` : `Yandex Browser`,
	`XI` : `Xiino`,

	// detected browsers in older versions
	// `IA` : `Iceape`,  : pim
	// `SM` : `SeaMonkey`,  : pim
}

// Browser families mapped to the short codes of the associated browsers
var browserFamilies = map[string][]string{
	`Android Browser`    : []string{`AN`, `MU`},
	`BlackBerry Browser` : []string{`BB`},
	`Baidu`              : []string{`BD`, `BS`},
	`Amiga`              : []string{`AV`, `AW`},
	`Chrome`             : []string{`CH`, `BA`, `BR`, `CC`, `CD`, `CM`, `CI`, `CF`, `CN`, `CR`, `CP`, `DD`, `IR`, `RM`, `AO`, `TS`, `VI`, `PT`, `AS`, `TB`, `AD`, `SB`, `WP`, `I3`, `CV`, `WH`, `SZ`, `QW`, `LF`, `KW`, `2B`, `CE`, `EC`, `MT`, `MS`, `HA`, `OC`},
	`Firefox`            : []string{`FF`, `FE`, `FM`, `SX`, `FB`, `PX`, `MB`, `EI`, `WF`, `CU`, `TF`, `QM`, `FR`, `I4`, `GZ`, `MO`, `F1`, `BI`},
	`Internet Explorer`  : []string{`IE`, `IM`, `PS`},
	`Konqueror`          : []string{`KO`},
	`NetFront`           : []string{`NF`},
	`NetSurf`            : []string{`NE`},
	`Nokia Browser`      : []string{`NB`, `NO`, `NV`, `DO`},
	`Opera`              : []string{`OP`, `OM`, `OI`, `ON`, `OO`, `OG`, `OH`, `O1`},
	`Safari`             : []string{`SF`, `MF`, `SO`},
	`Sailfish Browser`   : []string{`SA`},
}

// Browsers that are available for mobile devices only
var mobileOnlyBrowsers = []string{
	`36`, `OC`, `PU`, `SK`, `MF`, `OI`, `OM`, `DD`, `DB`, `ST`, `BL`, `IV`, `FM`, `C1`, `AL`, `SA`, `SB`, `FR`, `WP`, `HA`, `NX`, `HU`, `VV`, `RE`,
}

func GetBrowserFamily(browserLabel string) (string, bool) {
	for k, vs := range browserFamilies {
		for _, v := range vs {
			if v == browserLabel {
				return k, true
			}
		}
	}
	return "", false
}

// Returns if the given browser is mobile only
func IsMobileOnlyBrowser(browser string) bool {
	if ArrayContainsString(mobileOnlyBrowsers, browser) {
		return true
	}
	if v, ok := availableBrowsers[browser]; ok {
		return ArrayContainsString(mobileOnlyBrowsers, v)
	}
	return false
}

type BrowserMatchResult = ClientMatchResult

type Engine struct {
	Default  string            `yaml:"default" json:"default"`
	Versions map[string]string `yaml:"versions" json:"versions"`
}

type BrowserItem struct {
	Regular `yaml:",inline" json:",inline"`
	Name    string  `yaml:"name" json:"name"`
	Version string  `yaml:"version" json:"version"`
	Engine  *Engine `yaml:"engine" json:"engine"`
}

//  Client parser for browser detection
type Browser struct {
	Regexes      []*BrowserItem
	engine       BrowserEngine
	verCache     map[string]*Version
}

const ParserNameBrowser = `browser`
const FixtureFileBrowser = `browsers.yml`

func init() {
	RegClientParser(ParserNameBrowser,
		func(dir string) ClientParser {
			return NewBrowser(filepath.Join(dir, FixtureFileBrowser))
		})
}

func NewBrowser(fileName string) *Browser {
	b := &Browser{}
	b.engine.ParserName = ParserNameBrowserEngine
	if err := b.Load(fileName); err != nil {
		return nil
	}
	return b
}

func (b *Browser) Load(file string) error {
	b.verCache = make(map[string]*Version)
	var v []*BrowserItem
	err := ReadYamlFile(file, &v)
	if err != nil {
		return err
	}
	engineFile := file[0:len(file)-len(FixtureFileBrowser)] + FixtureFileBrowserEngine
	err = b.engine.Load(engineFile)
	if err != nil {
		return err
	}
	b.Regexes = v
	return nil
}

func (b *Browser) PreMatch(ua string) bool {
	return true
}

func (b *Browser) Parse(ua string) *BrowserMatchResult {
	for _, regex := range b.Regexes {
		matches := regex.MatchUserAgent(ua)
		if len(matches) > 0{
			name := BuildByMatch(regex.Name, matches)
			for browserShort, browserName := range availableBrowsers {
				if StringEqualIgnoreCase(name, browserName) {
					version := BuildVersion(regex.Version, matches)
					engine := b.BuildEngine(regex.Engine, version, ua)
					engineVersion := b.BuildEngineVersion(engine, ua)
					return &BrowserMatchResult{
						Type:          ParserNameBrowser,
						Name:          browserName,
						ShortName:     browserShort,
						Version:       version,
						Engine:        engine,
						EngineVersion: engineVersion,
					}
				}
			}
		}
	}
	return nil
}

func (b *Browser) BuildEngine(engineData *Engine, browserVersion, ua string) string {
	engine := ""
	if engineData != nil {
		engine = engineData.Default
		for version, versionEngine := range engineData.Versions {
			if gover.CompareSimple(browserVersion, version) >= 0 {
				engine = versionEngine
			}
		}
	}
	if engine == "" {
		if engineResult := b.engine.Parse(ua); engineResult != nil {
			engine = engineResult.Name
		}
	}
	return engine
}

func (b *Browser) BuildEngineVersion(engine, ua string) string {
	if engine == "" {
		return ""
	}
	v, ok := b.verCache[engine]
	if !ok {
		v = &Version{Engine: engine}
		v.Compile()
		b.verCache[engine] = v
	}
	return v.Parse(ua)
}
