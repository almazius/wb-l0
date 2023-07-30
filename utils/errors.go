package utils

import (
	"strconv"
	"strings"
)

// GetCodeFromMyError get error code from string
func GetCodeFromMyError(err error) int {
	words := strings.Split(err.Error(), " ")
	if len(words) > 1 {
		val, err := strconv.Atoi(words[1])
		if err == nil {
			return val
		}
	}
	return 0
}
