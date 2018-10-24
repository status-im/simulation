package main

import "net/http"

func allowCORS(fn func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Access-Control-Allow-Origin"); origin == "" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if origin := r.Header.Get("Origin"); origin != "" {
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				// set preflight options
				return
			}
		}

		fn(w, r)
	}
}
