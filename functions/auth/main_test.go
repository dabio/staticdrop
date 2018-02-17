package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/dabio/staticdrop/oauth2"
)

func TestHandler_noCode(t *testing.T) {
	t.Parallel()

	config := &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_APP_KEY"),
		ClientSecret: os.Getenv("DROPBOX_APP_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	rec := httptest.NewRecorder()
	h := handler(config)
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

	equals(t, 303, rec.Code)
}

func TestHandler_code(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"access_token": "ABCDEFG", "token_type": "bearer", "account_id": "dbid:AAH4f99T0taONIb-OurWxbNQ6ywGRopQngc", "uid": "12345"}`))
		}),
	)
	defer ts.Close()

	config := &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_APP_KEY"),
		ClientSecret: os.Getenv("DROPBOX_APP_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  ts.URL,
			TokenURL: ts.URL,
		},
	}

	rec := httptest.NewRecorder()
	h := handler(config)
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/?code=1234", nil))

	equals(t, 200, rec.Code)
}

func TestHandler_invalidCode(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
			w.Write([]byte(`{"error_description": "code doesn't exist or has expired", "error": "invalid_grant"}`))
		}),
	)
	defer ts.Close()

	config := &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_APP_KEY"),
		ClientSecret: os.Getenv("DROPBOX_APP_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  ts.URL,
			TokenURL: ts.URL,
		},
	}

	rec := httptest.NewRecorder()
	h := handler(config)
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/?code=1234", nil))

	equals(t, 200, rec.Code)
}

func TestFixRequestURL(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("GET", "/target", nil)
	fixRequestURL(req)
	equals(t, "https://example.com/target", req.URL.String())

	req = httptest.NewRequest("GET", "https://1234567890.execute-api.us-east-1.amazonaws.com/pets", nil)
	// mimic stage environment
	req.Header.Set("X-Stage", "prod")
	fixRequestURL(req)
	equals(t, "https://1234567890.execute-api.us-east-1.amazonaws.com/prod/pets", req.URL.String())
}

// equals fails the test if exp is not equal to act.
func equals(t testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}
