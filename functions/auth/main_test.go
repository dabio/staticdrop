package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestGetNotFound(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	handle(rec, httptest.NewRequest("GET", "/blah", nil))
	equals(t, 404, rec.Code)

	handle(rec, httptest.NewRequest("POST", "/", nil))
	equals(t, 404, rec.Code)

	handle(rec, httptest.NewRequest("DELETE", "/", nil))
	equals(t, 404, rec.Code)
}

func TestGet(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte(`{"access_token": "ABCDEFG", "token_type": "bearer", "account_id": "dbid:AAH4f99T0taONIb-OurWxbNQ6ywGRopQngc", "uid": "12345"}`))
		}),
	)
	defer srv.Close()

	rec := httptest.NewRecorder()

	tokenAPI = srv.URL
	handle(rec, httptest.NewRequest("GET", "/?"+codeParam+"=123", nil))

	equals(t, 200, rec.Code)
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
