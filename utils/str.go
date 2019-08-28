package utils

import (
	"fmt"
	"strings"
)

// GetPathLastSegment path separated by '/'
func GetPathLastSegment(path string) string {
	return LastItem(path, "/")
}

// StringConcatSlash concat `str`s with '/'
func StringConcatSlash(str ...string) string {
	return strings.Join(str, "/")
}

// LastItem of `str` separated by `sep`
func LastItem(str, sep string) string {
	parts := strings.Split(str, sep)
	return parts[len(parts)-1]
}

// FirstItem of `str` separated by `sep`
func FirstItem(str, sep string) string {
	item, _ := NthItem(str, sep, 0)
	return item
}

// NthItem of `str` separated by `sep`
func NthItem(str, sep string, n int) (string, error) {
	parts := strings.SplitN(str, sep, n+2)
	if len(parts) <= n {
		return "", fmt.Errorf("out of index")
	}
	return parts[n], nil
}
