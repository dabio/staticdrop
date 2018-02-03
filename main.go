package main

import (
	"log"
	"net/http"
	"os"
)

const query = "challenge"

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGET(w, r)
	case http.MethodPost:
		handlePOST(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleGET(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Query().Get(query)))
}

func handlePOST(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	// we need to put this into a queue here
}
