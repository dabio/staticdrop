package dropbox

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const dropboxTokenURL = "https://api.dropboxapi.com/oauth2/token"

// Auth contains all necessities to obtain the token.
type Auth struct {
	Code         string
	ClientID     string
	ClientSecret string
	URL          string
}

// Token contains the reponse token from the API.
type Token struct {
	AccessToken string
	TokenType   string
	AccountID   string
}

// Access calls the token API to acquire a bearer token.
func Access(auth *Auth) (*Token, error) {
	url := auth.URL
	if url == "" {
		url = dropboxTokenURL
	}

	url = fmt.Sprintf("%s?code=%s", url, auth.Code)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(auth.ClientID, auth.ClientSecret)

	client := &http.Client{
		Timeout: defaultTimeout,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var token *Token
	if err = json.NewDecoder(res.Body).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}
