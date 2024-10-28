package fidl

import "net/http"

func SetHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
}
