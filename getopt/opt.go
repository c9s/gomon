package getopt

import "strings"

type Option struct {
	LongFlag  string
	ShortFlag string
	Usage     string
	Value     interface{}
}

type Options *[]Option

func Opt(spec string, defvalue interface{}, usage string) {

}
