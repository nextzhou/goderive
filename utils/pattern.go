package utils

import "regexp"

var identPattern = regexp.MustCompile(`^[_a-zA-Z]\w*$`)

func ValidateIdentName(name string) bool {
	return identPattern.MatchString(name)
}
