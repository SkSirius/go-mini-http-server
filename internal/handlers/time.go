package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func TimeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSONResponse{Message: time.Now().Format(time.RFC3339)})
}
