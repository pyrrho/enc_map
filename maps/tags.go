package maps

import (
	"strings"
)

type tagOptions map[string]string

func parseTag(tag string) (string, tagOptions) {
	strs := strings.Split(tag, ",")
	name := strs[0]

	opts := make(tagOptions, len(strs)-1)
	for _, str := range strs[1:] {
		idx := strings.Index(str, "=")
		if idx < 0 {
			opts.setOption(str, "")
		} else {
			opts.setOption(str[:idx], str[:idx])
		}
	}
	return name, opts
}

func (opts tagOptions) setOption(option string, value string) {
	opts[strings.ToLower(option)] = value
}

func (opts tagOptions) getOption(option string) (val string, ok bool) {
	val, ok = opts[strings.ToLower(option)]
	return
}

func (opts tagOptions) Contains(option string) bool {
	_, ok := opts.getOption(option)
	return ok
}

func (opts tagOptions) ValueOf(option string) string {
	val, ok := opts.getOption(option)
	if !ok {
		return ""
	} else if val == "" {
		return option
	}
	return val
}
