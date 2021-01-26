package devicedetector

import (
	"strconv"
	"sync"
	"testing"

	"gotest.tools/assert"

	regexp "github.com/dlclark/regexp2"
	. "github.com/muxinc/devicedetector/parser"
	"github.com/muxinc/devicedetector/parser/client"
	"github.com/muxinc/devicedetector/parser/device"
)

var dd, _ = NewDeviceDetector("regexes")

func TestParseInvalidUA(t *testing.T) {
	info := dd.Parse(`12345`)
	if info != nil {
		t.Fatal("testParseInvalidUA fail")
	}
}

func TestInstanceReusage(t *testing.T) {
	userAgents := [][]string{
		[]string{
			`Sraf/3.0 (Linux i686 ; U; HbbTV/1.1.1 (+PVR+DL;NEXtUS; TV44; sw1.0) CE-HTML/1.0 Config(L:eng,CC:DEU); en/de)`,
			``,
			``,
		},
		[]string{
			`Mozilla/5.0 (Linux; Android 4.2.2; ARCHOS 101 PLATINUM Build/JDQ39) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.114 Safari/537.36`,
			`Archos`,
			`101 PLATINUM`,
		},
		[]string{
			`Opera/9.80 (Linux mips; U; HbbTV/1.1.1 (; Vestel; MB95; 1.0; 1.0; ); en) Presto/2.10.287 Version/12.00`,
			`Vestel`,
			`MB95`,
		},
	}

	for i, item := range userAgents {
		info := dd.Parse(item[0])
		assert.Equal(t, info.GetBrandName(), item[1], i)
		assert.Equal(t, info.GetModel(), item[2], i)
	}
}

func TestVersionTruncation(t *testing.T) {
	data := map[int][]string{
		VERSION_TRUNCATION_NONE:  []string{`Mozilla/5.0 (Linux; Android 4.2.2; ARCHOS 101 PLATINUM Build/JDQ39) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.114 Safari/537.36`, `4.2.2`, `34.0.1847.114`},
		VERSION_TRUNCATION_BUILD: []string{`Mozilla/5.0 (Linux; Android 4.2.2; ARCHOS 101 PLATINUM Build/JDQ39) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.114 Safari/537.36`, `4.2.2`, `34.0.1847.114`},
		VERSION_TRUNCATION_PATCH: []string{`Mozilla/5.0 (Linux; Android 4.2.2; ARCHOS 101 PLATINUM Build/JDQ39) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.114 Safari/537.36`, `4.2.2`, `34.0.1847`},
		VERSION_TRUNCATION_MINOR: []string{`Mozilla/5.0 (Linux; Android 4.2.2; ARCHOS 101 PLATINUM Build/JDQ39) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.114 Safari/537.36`, `4.2`, `34.0`},
		VERSION_TRUNCATION_MAJOR: []string{`Mozilla/5.0 (Linux; Android 4.2.2; ARCHOS 101 PLATINUM Build/JDQ39) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.114 Safari/537.36`, `4`, `34`},
	}
	for k, v := range data {
		SetVersionTruncation(k)
		info := dd.Parse(v[0])
		assert.Equal(t, info.GetOs().Version, v[1])
		assert.Equal(t, info.GetClient().Version, v[2])
	}
}

func TestBot(t *testing.T) {
	type BotTest struct {
		Ua  string         `yaml:"user_agent"`
		Bot BotMatchResult `yaml:"bot"`
	}
	var listBotTest []BotTest
	err := ReadYamlFile(`fixtures/bots.yml`, &listBotTest)
	if err != nil {
		t.Error(err)
	}

	for _, item := range listBotTest {
		info := dd.Parse(item.Ua)
		bot := info.GetBot()
		if bot == nil {
			t.Error("bot is null")
		}

		if item.Bot.Equal(bot) == false {
			t.Error("bot is null")
		}

		osName := info.GetOs().ShortName
		clientName := info.GetClient().ShortName
		assert.Equal(t, osName, "", osName)
		assert.Equal(t, clientName, "", clientName)
	}
}

func TestTypeMethods(t *testing.T) {
	data := map[string][]bool{
		`Googlebot/2.1 (http://www.googlebot.com/bot.html)`: []bool{true, false, false},
		`Mozilla/5.0 (Linux; Android 4.4.2; Nexus 4 Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.136 Mobile Safari/537.36`: []bool{false, true, false},
		`Mozilla/5.0 (Linux; Android 4.4.3; Build/KTU84L) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.117 Mobile Safari/537.36`:         []bool{false, true, false},
		`Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)`:                                                                    []bool{false, false, true},
		`Mozilla/3.01 (compatible;)`: []bool{false, false, false},
		// Mobile only browsers:
		`Opera/9.80 (J2ME/MIDP; Opera Mini/9.5/37.8069; U; en) Presto/2.12.423 Version/12.16`:                                                                                          []bool{false, true, false},
		`Mozilla/5.0 (X11; U; Linux i686; th-TH@calendar=gregorian) AppleWebKit/534.12 (KHTML, like Gecko) Puffin/1.3.2665MS Safari/534.12`:                                            []bool{false, true, false},
		`Mozilla/5.0 (Linux; Android 4.4.4; MX4 Pro Build/KTU84P) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/33.0.0.0 Mobile Safari/537.36; 360 Aphone Browser (6.9.7)`: []bool{false, true, false},
		`Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_5_7; xx) AppleWebKit/530.17 (KHTML, like Gecko) Version/4.0 Safari/530.17 Skyfire/6DE`:                                           []bool{false, true, false},
		// useragent containing non unicode chars
		`Mozilla/5.0 (Linux; U; Android 4.1.2; ru-ru; PMP7380D3G Build/JZO54K) AppleWebKit/534.30 (KHTML, ÃÂºÃÂ°ÃÂº Gecko) Version/4.0 Safari/534.30`: []bool{false, true, false},
	}
	for k, v := range data {
		dd.DiscardBotInformation = true
		info := dd.Parse(k)
		assert.Equal(t, info.IsBot(), v[0], k)
		assert.Equal(t, info.IsMobile(), v[1], k)
		assert.Equal(t, info.IsDesktop(), v[2], k)
	}
}

func TestGetOs(t *testing.T) {
	ua := `Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)`
	info := dd.Parse(ua)
	os := info.GetOs()
	assert.Equal(t, os.Name, `Windows`)
	assert.Equal(t, os.ShortName, `WIN`)
	assert.Equal(t, os.Version, `7`)
	assert.Equal(t, os.Platform, `x64`)
}

func TestGetClient(t *testing.T) {
	ua := `Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)`
	info := dd.Parse(ua)
	client := info.GetClient()
	assert.Equal(t, client.Type, `browser`)
	assert.Equal(t, client.Name, `Internet Explorer`)
	assert.Equal(t, client.ShortName, `IE`)
	assert.Equal(t, client.Version, `9.0`)
	assert.Equal(t, client.Engine, `Trident`)
	assert.Equal(t, client.EngineVersion, `5.0`)
}

func TestGetBrandName(t *testing.T) {
	ua := `Mozilla/5.0 (Linux; Android 4.4.2; Nexus 4 Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/33.0.1750.136 Mobile Safari/537.36`
	info := dd.Parse(ua)
	assert.Equal(t, info.GetBrandName(), `Google`)
}

func TestIsTouchEnabled(t *testing.T) {
	ua := `Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; ARM; Trident/6.0; Touch; ARMBJS)`
	info := dd.Parse(ua)
	assert.Check(t, info.IsTouchEnabled())
}

func TestSkipBotDetection(t *testing.T) {
	ua := `Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)`
	info := dd.Parse(ua)
	assert.Check(t, info.IsMobile() == false)
	assert.Check(t, info.IsBot())
	dd.SkipBotDetection = true
	info = dd.Parse(ua)
	assert.Check(t, info.IsMobile())
	assert.Check(t, info.IsBot() == false)
}

type SmartFixture struct {
	UserAgent     string                    `yaml:"user_agent"`
	Os            *OsMatchResult            `yaml:"os"`
	Client        *client.ClientMatchResult `yaml:"client"`
	Device        *device.DeviceMatchResult `yaml:"device"`
	OsFamily      string                    `yaml:"os_family"`
	BrowserFamily string                    `yaml:"browser_family"`
}

func TestRegThread(t *testing.T) {
	// read file
	var lists [][]*SmartFixture
	for i := 0; i <= 21; i++ {
		var list []*SmartFixture
		var name string
		if i == 0 {
			name = `smartphone.yml`
		} else {
			name = `smartphone-` + strconv.Itoa(i) + `.yml`
		}
		err := ReadYamlFile(`fixtures/`+name, &list)
		if err == nil {
			lists = append(lists, list)
		}
	}
	rs := []*regexp.Regexp{adrMobReg, touchReg, adrTabReg, chrMobReg, chrTabReg, opaTabReg, opaTvReg}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		for _, list := range lists {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for _, f := range list {
					ua := f.UserAgent
					for _, reg := range rs {
						reg.MatchString(ua)
						reg.FindStringMatch(ua)
					}
				}
			}()
		}
	}
	wg.Wait()
}
