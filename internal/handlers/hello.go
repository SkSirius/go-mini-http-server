package handlers

import (
	"encoding/json"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JSONResponse{Message: "Hello, world!"})
}
