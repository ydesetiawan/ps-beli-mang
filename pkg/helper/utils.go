package helper

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidateISO8601(fl validator.FieldLevel) bool {
	// List of layouts for various ISO 8601 formats
	dateStr, ok := fl.Field().Interface().(string)
	if !ok {
		// Field is not a string
		return false
	}
	layouts := []string{
		time.RFC3339,                          // Full date and time
		"2006-01-02",                          // Date only
		"2006-01-02T15:04:05Z07:00",           // Date and time with time zone
		"2006-01-02T15:04:05",                 // Date and time without time zone
		"2006-01-02T15:04Z07:00",              // Date and time with time zone, no seconds
		"2006-01-02T15:04",                    // Date and time without time zone, no seconds
		"2006-01-02T15:04:05.999999999Z07:00", // Date and time with nanoseconds
	}

	// Try parsing the input string with each layout
	for _, layout := range layouts {
		if _, err := time.Parse(layout, dateStr); err == nil {
			return true
		}
	}
	return false
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

func ValidateNIPForIT(fl validator.FieldLevel) bool {
	nipInt := fl.Field().Int()
	nip := strconv.Itoa(int(nipInt))

	return ValidatePrefixIT(nip)
}

func ValidatePrefixIT(nip string) bool {
	// Check the length
	if len(nip) < 13 || len(nip) > 15 {
		return false
	}

	// Check the prefix "615"
	if nip[:3] != "615" {
		return false
	}

	return true
}

func ValidatePrefixNurse(nip string) bool {
	// Check the length
	if len(nip) < 13 || len(nip) > 15 {
		return false
	}

	// Check the prefix "303"
	if nip[:3] != "303" {
		return false
	}

	return true
}

func ValidateNIPForNurse(fl validator.FieldLevel) bool {
	nipInt := fl.Field().Int()
	nip := strconv.Itoa(int(nipInt))

	return ValidatePrefixNurse(nip)
}

// ValidateNIP checks if the NIP meets the specific criteria
func ValidateNIP(fl validator.FieldLevel) bool {
	nipInt := fl.Field().Int()
	nip := strconv.Itoa(int(nipInt))

	return ValidateNIPFromStr(nip)
}

func ValidateNIPFromStr(nip string) bool {
	genderDigit := nip[3]
	if genderDigit != '1' && genderDigit != '2' {
		return false
	}

	yearStr := nip[4:8]
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 || year > (time.Now().Year()) {
		return false
	}

	monthStr := nip[8:10]
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return false
	}

	randomDigits := nip[10:]
	if len(randomDigits) < 3 || len(randomDigits) > 5 {
		return false
	}
	if _, err := strconv.Atoi(randomDigits); err != nil {
		return false
	}

	return true
}

func GenerateULID() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())

	return ulid.MustNew(ms, entropy).String()

}
