package helper

import "encoding/json"

func JsonEncodeDefaultEmpty(v any) string {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
