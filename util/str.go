package util

import "strings"

func ContainNumber(s string) bool {
	return strings.ContainsAny(s, "0123456789")
}
