package getopt

import "strings"

type Option struct {
	LongFlag  string
	ShortFlag string
	Usage     string
	Value     interface{}
}

type Options *[]Option

func Opt(spec string, val interface{}, usage string) *Option {
	opt := Option{}
	parts := strings.SplitN(spec, "|", 2)

	// with short and long flag
	if len(parts) == 1 {
		if len(parts[0]) == 1 {
			opt.ShortFlag = parts[0]
		} else {
			opt.LongFlag = parts[0]
		}
		opt.Usage = usage
		opt.Value = val
	} else if len(parts) == 2 {
		opt.ShortFlag = parts[0]
		opt.LongFlag = parts[1]
		opt.Usage = usage
		opt.Value = val
	}
	return &opt
}
