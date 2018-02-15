package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dabio/staticdrop/oauth2"
)

const (
	codeParam = "code"
	localhost = "localhost"
	tokenURL  = "https://api.dropboxapi.com/oauth2/token"
	authURL   = "https://www.dropbox.com/oauth2/authorize"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if !strings.HasPrefix(r.Host, localhost) {
		scheme = scheme + "s"
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_API_KEY"),
		ClientSecret: os.Getenv("DROPBOX_API_SECRET"),
		RedirectURL:  fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.String()),
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	code := r.URL.Query().Get(codeParam)
	if code == "" {
		log.Print(`{"message":"redirect to authentication page"}`)
		http.Redirect(w, r, config.AuthCodeURL(), http.StatusSeeOther)
		return
	}
	w.Write([]byte(config.AuthCodeURL()))
}
