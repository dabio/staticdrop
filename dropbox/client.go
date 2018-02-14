package dropbox

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// Client implements the dropbox client used to get information from the
// dropbox API.
type Client struct {
	*Config

	Files *Files
}

// New creates a new client for quering the API.
func New(config *Config) *Client {
	c := &Client{Config: config}
	c.Files = &Files{c}

	return c
}

func (c *Client) call(path string, in interface{}) (io.ReadCloser, error) {
	url := c.URL + path

	body, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	rec, _, err := c.do(req)

	return rec, err
}

func (c *Client) do(req *http.Request) (io.ReadCloser, int64, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	if res.StatusCode < 400 {
		return res.Body, res.ContentLength, err
	}
	defer res.Body.Close()

	e := &Error{
		Status:     http.StatusText(res.StatusCode),
		StatusCode: res.StatusCode,
	}

	if err := json.NewDecoder(res.Body).Decode(e); err != nil {
		return nil, 0, err
	}

	return nil, 0, e
}
