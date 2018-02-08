package main

import (
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()

	var url string
	challenges := []string{
		"sth",
		"very",
		"random",
		"7c3ab00c-4fd9-413c-8d99-603ecfcf2e1d",
	}

	for _, challenge := range challenges {
		url = fmt.Sprintf("/?%s=%s", challengeQuery, challenge)

		rec := httptest.NewRecorder()
		handle(rec, httptest.NewRequest("GET", url, nil))

		equals(t, 200, rec.Code)
		equals(t, challenge, rec.Body.String())
	}
}

func TestPostNothing(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	handle(rec, httptest.NewRequest("POST", "/", nil))

	equals(t, 403, rec.Code)
}

func TestPostSignature(t *testing.T) {
	t.Parallel()

	body := `{"list_folder": {"accounts": ["dbid:AADuYvHN_tjvjM4JyjA70vGwIMNau360Mbo"]}, "delta": {"users": [787701]}}`
	signature := "a2f62303570360415ffb6cd3ddc2e39500b59ee0ae684df0a1de8630d555d2d0"

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Add(signatureHeader, signature)
	handle(rec, req)

	equals(t, 200, rec.Code)
}

func Test404(t *testing.T) {
	t.Parallel()

	paths := []string{"/new", "/path", "/different/path"}

	for _, path := range paths {
		rec := httptest.NewRecorder()
		handle(rec, httptest.NewRequest("GET", path, nil))

		equals(t, 404, rec.Code)
	}
}

func TestNotAllowed(t *testing.T) {
	t.Parallel()

	methods := []string{"PUT", "PATCH", "DELETE", "OPTIONS"}

	for _, method := range methods {
		rec := httptest.NewRecorder()
		handle(rec, httptest.NewRequest(method, "/", nil))

		equals(t, 405, rec.Code)
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
