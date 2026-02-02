package res

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResDTO struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

func Json(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println("Encode error", err)
		return
	}
}
