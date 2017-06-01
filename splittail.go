package golax

import "strings"

func SplitTail(s, sep string) []string {
	parts := strings.Split(s, sep)
	l := len(parts)
	if 1 == l {
		return []string{s}
	}
	return []string{strings.Join(parts[:l-1], sep), parts[l-1]}
}
