package oauth2_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/dabio/staticdrop/oauth2"
)

const (
	clientID     = "CLIENT_ID"
	clientSecret = "CLIENT_SECRET"
	redirectURL  = "REDIRECT_URL"
)

func config(url string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  url,
			TokenURL: url,
		},
	}
}

func TestAuthCodeURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		url string
		exp string
	}{
		{url: "url", exp: fmt.Sprintf("url?client_id=%s&redirect_uri=%s&response_type=code", clientID, redirectURL)},
		{url: "url?blah", exp: fmt.Sprintf("url?blah&client_id=%s&redirect_uri=%s&response_type=code", clientID, redirectURL)},
	}

	for _, test := range tests {
		c := config(test.url)
		equals(t, test.exp, c.AuthCodeURL())

	}
}

func TestExchangeCode2Token(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			equals(t, "POST", r.Method)
			equals(t, r.Header.Get("Content-Type"), "application/x-www-form-urlencoded")

			w.Write([]byte(`{"access_token": "ABCDEFG", "token_type": "bearer", "account_id": "dbid:AAH4f99T0taONIb-OurWxbNQ6ywGRopQngc", "uid": "12345"}`))
		}),
	)
	defer ts.Close()

	config := config(ts.URL)
	token, err := config.Exchange("123456")

	ok(t, err)
	equals(t, &oauth2.Token{AccessToken: "ABCDEFG", TokenType: "bearer"}, token)
	equals(t, "bearer", token.Type())
}

func TestExchangeCodeRetrieveError(t *testing.T) {
	t.Parallel()

	const json = `{"error_description": "code doesn't exist or has expired", "error": "invalid_grant"}`
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			equals(t, "POST", r.Method)
			equals(t, r.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
			w.WriteHeader(400)
			w.Write([]byte(json))
		}),
	)
	defer ts.Close()

	config := config(ts.URL)
	_, err := config.Exchange("123456")

	if err == nil {
		t.Error("want RetrieveError")
	}
}

// equals fails the test if exp is not equal to act
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}
