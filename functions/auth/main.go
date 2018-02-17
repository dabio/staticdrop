package main

import (
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
	// awsSuffix is the ending hostname when the app runs using the api gateway
	// hostname.
	awsSuffix = ".amazonaws.com"
)

func main() {
	config := &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_APP_KEY"),
		ClientSecret: os.Getenv("DROPBOX_APP_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	http.Handle("/", handler(config))

	addr := ":" + os.Getenv("PORT")
	log.Fatal(gateway.ListenAndServe(addr, nil))
}

func handler(config *oauth2.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fixRequestURL(r)
		config.RedirectURL = r.URL.String()

		code := r.URL.Query().Get(codeParam)
		if code == "" {
			log.Print(`{"message":"redirect to authentication page"}`)
			http.Redirect(w, r, config.AuthCodeURL(), http.StatusSeeOther)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Printf(`{"message":"exchange code", "code": "%s"}`, code)

		// defer func() {
		token, err := config.Exchange(code)
		if err != nil {
			log.Printf(`{"message":"error getting token", "error": "%s"}`, err)
		} else {
			log.Printf(`{"message":"got the token", "token": "%s", "type": "%s"}`, token.AccessToken, token.TokenType)
		}
		// }()
	})
}

func fixRequestURL(r *http.Request) {
	if strings.HasSuffix(r.Host, awsSuffix) {
		r.URL.Path = "/" + r.Header.Get("X-Stage") + r.URL.Path
	}
	if r.Host != localhost && r.URL.Scheme == "" {
		r.URL.Scheme = "https"
	}
	if r.URL.Host == "" {
		r.URL.Host = r.Host
	}
}
