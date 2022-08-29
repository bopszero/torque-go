package utils

import "strings"

func IsSameStringCI(left, right string) bool {
	return strings.ToLower(left) == strings.ToLower(right)
}

func StringTrim(text string) string {
	return strings.Trim(text, " \t")
}
