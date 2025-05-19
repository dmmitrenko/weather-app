package http

import (
	_ "embed"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed static/index.html
var indexHTML []byte

func RegisterStatic(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	}).Methods("GET")
}
