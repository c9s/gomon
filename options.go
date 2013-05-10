package main

type gomonOption struct {
	flag        string
	value       interface{}
	description string
}

type gomonOptions []*gomonOption

var options = gomonOptions{
	{"h", false, "Show Help"},
	{"b", false, "Run `go build`, the default behavior"},
	{"t", false, "Run `go test`"},
	{"i", false, "Run `go install`"},
	{"f", false, "Run `go fmt`"},
	{"r", false, "Run `go run`"},
	{"x", false, "Show verbose command"},
	{"v", false, "Show version"},
	{"d", false, "Print debug message"},
	{"chdir", false, "Run commands on directory"},
	{"growl", false, "Use Growler"},
	{"install-growl-icons", false, "Install growl icons"},
	{"gntp", "127.0.0.1:23053", "The GNTP DSN"},
}

func (options gomonOptions) Get(flag string) *gomonOption {
	for _, option := range options {
		if option.flag == flag {
			return option
		}
	}
	return nil
}

func (options gomonOptions) String(flag string) string {
	for _, option := range options {
		if option.flag == flag {
			s, _ := option.value.(string)
			return s
		}
	}
	return ""
}

func (options gomonOptions) Bool(flag string) bool {
	for _, option := range options {
		if option.flag == flag {
			b, _ := option.value.(bool)
			return b
		}
	}
	return false
}

func (options gomonOptions) IsBool(flag string) bool {
	for _, option := range options {
		if option.flag == flag {
			_, ok := option.value.(bool)
			return ok
		}
	}
	return false
}
