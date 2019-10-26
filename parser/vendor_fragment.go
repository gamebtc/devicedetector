package parser

const ParserNameVendor = "vendorfragments"
const FixtureFileVendor = "vendorfragments.yml"

// Device parser for vendor fragment detection
type VendorFragments struct {
	vendorRegexes map[string][]*Regular
}

func NewVendor(file string) (*VendorFragments, error) {
	var m map[string][]string
	err := ReadYamlFile(file, &m)
	if err != nil {
		return nil, err
	}

	vendorRegexes := make(map[string][]*Regular)
	for name, brands := range m {
		for _, brand := range brands {
			regex := &Regular{Regex: brand + "[^a-z0-9]+"}
			regex.Compile()
			if regexes, ok := vendorRegexes[name]; ok {
				vendorRegexes[name] = append(regexes, regex)
			} else {
				regexes = make([]*Regular, 0, len(brands))
				vendorRegexes[name] = append(regexes, regex)
			}
		}
	}

	return &VendorFragments{
		vendorRegexes: vendorRegexes,
	}, nil
}

func (v *VendorFragments) Parse(ua string) string {
	for brand, regexes := range v.vendorRegexes {
		for _, regex := range regexes {
			if regex.IsMatchUserAgent(ua) {
				return GetShortName(brand)
			}
		}
	}
	return ""
}
