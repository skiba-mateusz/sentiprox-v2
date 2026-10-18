package respond

import (
	"encoding/json"
	"net/http"
)

type envelope struct {
	Data any `json:"data"`
}

type errorBody struct {
	Error string `json:"error"`
}

func write(w http.ResponseWriter, status int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}

func JSON(w http.ResponseWriter, status int, v any) error {
	return write(w, status, envelope{Data: v})
}

func Error(w http.ResponseWriter, status int, msg string) error {
	return write(w, status, errorBody{Error: msg})
}

func Decode(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
