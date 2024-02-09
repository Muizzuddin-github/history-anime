package response

import "net/http"

func SendJSONResponse(w http.ResponseWriter, status int, res []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}
