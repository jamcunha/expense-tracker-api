package api

import "net/http"

func NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, struct{}{})
	})

	return r
}
