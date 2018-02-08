package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const codeParam = "code"

// this is a variable so we can overwrite this in our test
var tokenAPI = "https://api.dropboxapi.com/oauth2/token"

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

	w.WriteHeader(http.StatusOK)
	log.Printf(`{"code":"%s"}`, code)

	// gather the bearer now
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s?code=%s&grant_type=authorization_code", tokenAPI, code),
		nil,
	)
	if err != nil {
		log.Printf(`{
			"message":"cannot connect gather bearer from dropbox",
			"error":"%s"
		}`, err)
		return
	}

	req.SetBasicAuth(os.Getenv("DROPBOX_API_KEY"), os.Getenv("DROPBOX_API_SECRET"))

	log.Print(tokenAPI)
	w.Write([]byte(tokenAPI))
	w.Write([]byte("https://www.dropbox.com/oauth2/authorize?client_id=" + os.Getenv("DROPBOX_APP_KEY") + "&response_type=code&redirect_uri=http://localhost:3000"))
}
