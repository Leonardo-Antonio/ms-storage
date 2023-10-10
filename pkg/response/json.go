package response

import (
	"encoding/json"
	"net/http"
	"time"
)

func Json(w http.ResponseWriter, data Response, statusCode int) {
	data.TimeStamp = uint64(time.Now().UnixMilli())
	jsonBytes, err := json.Marshal(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonBytes)
}
