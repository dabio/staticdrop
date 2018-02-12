package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dabio/staticdrop/oauth2"
)

const (
	codeParam = "code"
	tokenURL  = "https://api.dropboxapi.com/oauth2/token"
	authURL   = "https://www.dropbox.com/oauth2/authorize"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	code := r.URL.Query().Get(codeParam)
	if code == "" {
		log.Print(`{"message":"no code passed as query param"}`)
		http.NotFound(w, r)
		return
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_API_KEY"),
		ClientSecret: os.Getenv("DROPBOX_API_SECRET"),
		RedirectURL:  "http://localhost:3000",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
	w.Write([]byte(config.AuthCodeURL()))
}
