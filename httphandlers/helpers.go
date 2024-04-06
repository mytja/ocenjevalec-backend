package httphandlers

import (
	"encoding/json"
	"net/http"
)

func DumpJSON(jsonstruct interface{}) []byte {
	marshal, _ := json.Marshal(jsonstruct)
	return marshal
}

func WriteJSON(w http.ResponseWriter, jsonstruct interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(DumpJSON(jsonstruct))
}

func WriteForbiddenJWT(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "application/json")
	w.Write(DumpJSON(Response{Success: false, Data: "Forbidden"}))
}

func GetToken(r *http.Request) string {
	c, err := r.Cookie("Authorization")
	if err != nil {
		return ""
	}
	return c.Value
}

func WriteBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write(DumpJSON(Response{Success: false, Data: "Bad request"}))
}
