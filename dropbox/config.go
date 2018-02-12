package dropbox

import (
	"net/http"
	"time"
)

const (
	defaultTimeout = 5
	dropboxURL     = "https://api.dropboxapi.com/2"
)

// Config for the dropbox client
type Config struct {
	HTTPClient  *http.Client
	AccessToken string
	URL         string
}

// NewConfig return a Config with the given access token.
func NewConfig(accessToken string) *Config {
	return &Config{
		HTTPClient: &http.Client{
			Timeout: defaultTimeout * time.Second,
		},
		AccessToken: accessToken,
		URL:         dropboxURL,
	}
}
