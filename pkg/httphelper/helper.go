package httphelper

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(rw http.ResponseWriter, httpCode int, data interface{}) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	rw.Header().Set("Access-Control-Allow-Headers", "*")
	rw.Header().Set("Content-Type", "application/json")
	if httpCode != 0 {
		rw.WriteHeader(httpCode)
	}
	if data != nil {
		_ = json.NewEncoder(rw).Encode(data)
	}
}
