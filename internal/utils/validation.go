package utils

import "regexp"

var (
	lowerRegex = regexp.MustCompile(`[a-z]`)
	upperRegex = regexp.MustCompile(`[A-Z]`)
	digitRegex = regexp.MustCompile(`\d`)
)

func IsValidPassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}
	return lowerRegex.MatchString(pw) &&
		upperRegex.MatchString(pw) &&
		digitRegex.MatchString(pw)
}
