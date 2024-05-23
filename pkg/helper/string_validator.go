package helper

import (
	"golang.org/x/exp/slog"
	"strconv"
)

func ContainsEmoji(inp string) bool {
	if inp == "" {
		return false
	}
	foundEmoji := false
	res := ""
	runes := []rune(inp)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r < 128 {
			res += string(r)
		} else {
			res += "&#" + strconv.FormatInt(int64(r), 10) + ";"
			foundEmoji = true
		}
	}
	slog.Debug("Helper ContainsEmoji", "input", inp, "result", res)
	return foundEmoji
}
