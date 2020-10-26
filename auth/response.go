package auth

import (
	"encoding/json"
	"net/http"
)

func errorInResponse(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status) //HTTP status コードが入る
	json.NewEncoder(w).Encode(error)
	return
}

func responseByJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
	return
}
