package main

import (
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestRedirect(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	handle(rec, httptest.NewRequest("GET", "/", nil))

	equals(t, 303, rec.Code)
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
