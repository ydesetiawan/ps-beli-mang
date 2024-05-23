package helper

import (
	"strconv"
)

func ValidBool(value string) bool {
	_, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return true
}

func ContainsAll(superset, subset []string) bool {
	set := make(map[string]bool)
	for i := range superset {
		set[superset[i]] = true
	}
	for i := range subset {
		if !set[subset[i]] {
			return false
		}
	}
	return true
}
