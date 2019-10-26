package device

var deviceFactory = make(map[string]func(string) DeviceParser, 10)

func RegDeviceParser(name string, f func(string) DeviceParser) {
	deviceFactory[name] = f
}

func GetDeviceCreater(name string) func(string) DeviceParser {
	f, _ := deviceFactory[name]
	return f
}

func NewDeviceParser(dir, name string) DeviceParser {
	if f, ok := deviceFactory[name]; ok {
		return f(dir)
	}
	return nil
}

func NewDeviceParsers(dir string, names []string) []DeviceParser {
	r := make([]DeviceParser, len(names))
	for i, name := range names {
		if f, ok := deviceFactory[name]; ok {
			r[i] = f(dir)
		}
	}
	return r
}

