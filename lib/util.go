package lib

import (
	"regexp"
	"strings"
)

var (
	clean_re = regexp.MustCompile("[]{}|~`^[_?@/<>:;*=+()&\"'%#$.,!]")
)

func CleanName(name string) string {
	clean_name := clean_re.ReplaceAllString(name, "")
	clean_name = strings.ToLower(strings.Replace(clean_name, "-", " ", -1))
	return clean_name
}

func SplitName(name string) (name_map map[string]bool) {
	name_parts := strings.Split(name, " ")
	name_map = make(map[string]bool, len(name_parts))
	for _, part := range name_parts {
		name_map[part] = true
	}
	return
}

func MatchNames(name_map_a map[string]bool, name_map_b map[string]bool) bool {
	for key, _ := range name_map_b {
		_, found := name_map_a[key]
		if !found {
			return false
		}
	}
	return true
}
