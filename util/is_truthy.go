package util

import "strings"

var nos = [...]string{
	"",
	"0",
	"f",
	"false",
	"n",
	"no",
}

// IsTruthy returns `false` if the argument sounds like "false" (empty string,
// "0", "f", "false", and so on), and `true` otherwise.
func IsTruthy(val string) bool {
	val = strings.ToLower(val)

	for _, no := range nos {
		if val == no {
			return false
		}
	}

	return true
}
