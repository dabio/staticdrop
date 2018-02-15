package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/apex/gateway"
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
	log.Fatal(gateway.ListenAndServe(addr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if !strings.HasPrefix(r.Host, localhost) {
		scheme = scheme + "s"
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_APP_KEY"),
		ClientSecret: os.Getenv("DROPBOX_APP_SECRET"),
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

	w.WriteHeader(http.StatusOK)
	log.Printf(`{"message":"exchange code", "code": "%s"}`, code)

	defer func() {
		token, err := config.Exchange(code)
		if err != nil {
			log.Printf(`{"message":"error getting token", "error": "%s"}`, err)
		}

		log.Printf(`{"message":"got the token", "token": "%s", "type": "%s"}`, token.AccessToken, token.TokenType)
	}()
}
