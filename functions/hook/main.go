package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/apex/gateway"
)

const (
	challengeQuery  = "challenge"
	signatureHeader = "X-Dropbox-Signature"
)

var signKey = []byte(os.Getenv("DROPBOX_APP_SECRET"))

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handle)
	log.Fatal(gateway.ListenAndServe(addr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte(r.URL.Query().Get(challengeQuery)))
}

func handlePOST(w http.ResponseWriter, r *http.Request) {
	encoded := r.Header.Get(signatureHeader)
	signature, _ := hex.DecodeString(encoded)

	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if !checkHMAC(body, signature, signKey) {
		log.Printf(`{"message": "hmac failed", "body": "%s", "encoded": "%s", "signature": "%s"}`, body, encoded, signature)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	log.Printf(`{"body": "%s"}`, body)

	w.WriteHeader(http.StatusOK)
}

func checkHMAC(message, messageHMAC, key []byte) bool {
	h := hmac.New(sha256.New, key)
	h.Write(message)
	mac := h.Sum(nil)

	return hmac.Equal(messageHMAC, mac)
}
