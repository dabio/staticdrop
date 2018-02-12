package oauth2

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Endpoint contains the authorization and token URLs.
type Endpoint struct {
	AuthURL  string
	TokenURL string
}

// Config contains all informations needed for the oauth2 authorization flow.
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Endpoint     Endpoint
}

// AuthCodeURL returns the web page that lets the user sign in to the
// service and authorize your app.
func (c *Config) AuthCodeURL() string {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.AuthURL)
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {c.ClientID},
	}
	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}
	if strings.Contains(c.Endpoint.AuthURL, "?") {
		buf.WriteString("&")
	} else {
		buf.WriteString("?")
	}
	buf.WriteString(v.Encode())

	return buf.String()
}

// Exchange converts an authorization code into a token.
func (c *Config) Exchange(code string) (*Token, error) {
	v := url.Values{
		"grand_type": {"authorization_code"},
		"code":       {code},
	}
	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}

	return retrieveToken(c, v)
}

// Token contains the credentials used to authorize the requests to access
// protected resources on the OAuth 2.0 provider's backend.
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// Type returns the type of the token.
func (t *Token) Type() string {
	if strings.EqualFold(t.TokenType, "bearer") {
		return "bearer"
	}

	return "bearer"
}

// SetAuthHeader sets the authorization header of the given requests.
func (t *Token) SetAuthHeader(r *http.Request) {
	r.Header.Set("Authorization", t.Type()+" "+t.AccessToken)
}

func retrieveToken(c *Config, v url.Values) (token *Token, err error) {
	req, err := http.NewRequest(
		http.MethodPost,
		c.Endpoint.TokenURL,
		strings.NewReader(v.Encode()),
	)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&token)
	return
}
