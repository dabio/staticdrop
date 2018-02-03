package main

import (
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
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
		url = fmt.Sprintf("/?%s=%s", query, challenge)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		handle(rec, req)

		equals(t, 200, rec.Code)
		equals(t, challenge, rec.Body.String())
	}
}

func TestPost(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", nil)
	handle(rec, req)

	equals(t, 200, rec.Code)
}

func Test404(t *testing.T) {
	t.Parallel()

	paths := []string{"/new", "/path", "/different/path"}

	for _, path := range paths {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", path, nil)
		handle(rec, req)

		equals(t, 404, rec.Code)
	}
}

func TestNotAllowed(t *testing.T) {
	t.Parallel()

	methods := []string{"PUT", "PATCH", "DELETE", "OPTIONS"}

	for _, method := range methods {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/", nil)
		handle(rec, req)

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
