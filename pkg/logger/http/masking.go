package http

import "regexp"

// MaskingFunc is a function type to mask sensitive data inside the request and response body.
// it accepts a string of request/response body, and returns a masked string of it.
type MaskingFunc func(s string) string

var (
	sensitiveDataRegex = regexp.MustCompile(`"(password|client_secret|access_token|refresh_token)":"([^"]+)"`)
	replacerRegex      = regexp.MustCompile(`:"([^"]+)"`)
)

func defaultMaskingFunc(s string) string {
	return sensitiveDataRegex.ReplaceAllStringFunc(s, func(s string) string {
		return replacerRegex.ReplaceAllString(s, `:"${2}[FILTERED]"`)
	})
}
