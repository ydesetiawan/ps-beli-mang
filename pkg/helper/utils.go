package helper

import (
	"math/rand"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
)

func Keys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func IntKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func GenerateULID() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())

	return ulid.MustNew(ms, entropy).String()

}

func ValidateURL(fl validator.FieldLevel) bool {
	url, ok := fl.Field().Interface().(string)
	if !ok {
		// Field is not a string
		return false
	}
	// Define the regex pattern
	pattern := `^(http|https):\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(\/\S*)?$`
	// Match the regex pattern
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}
