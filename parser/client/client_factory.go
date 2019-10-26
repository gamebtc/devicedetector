package client

import "fmt"

var clientFactory = make(map[string]func(string) ClientParser, 10)

func RegClientParser(name string, f func(string) ClientParser) {
	clientFactory[name] = f
}

func GetClientCreater(name string) func(string) ClientParser {
	f, _ := clientFactory[name]
	return f
}

func NewClientParser(dir, name string) ClientParser {
	if f, ok := clientFactory[name]; ok {
		return f(dir)
	}
	return nil
}

func NewClientParsers(dir string, names []string) []ClientParser {
	r := make([]ClientParser, len(names))
	for i, name := range names {
		if f, ok := clientFactory[name]; ok {
			r[i] = f(dir)
		}
		if r[i] == nil {
			fmt.Printf("Client is null:" + name)
		}

	}
	return r
}

