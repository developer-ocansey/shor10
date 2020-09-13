package utils

import (
	"encoding/json"
	"net/http"
)

// FormatResponse ..
func FormatResponse(status string, message string, w http.ResponseWriter) {
	res := make(map[string]string)
	res["status"] = status
	res["message"] = message
	json.NewEncoder(w).Encode(res)
	return
}
