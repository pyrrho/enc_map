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
			opts[strings.ToLower(str)] = ""
		} else {
			opts[strings.ToLower(str[:idx])] = str[idx+1:]
		}
	}
	return name, opts
}

func (opts tagOptions) Contain(option string) bool {
	_, ok := opts[strings.ToLower(option)]
	return ok
}

func (opts tagOptions) ValueOf(option string) string {
	val, ok := opts[option]
	if !ok {
		return ""
	} else if val == "" {
		return option
	}
	return val
}
