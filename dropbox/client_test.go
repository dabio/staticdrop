package dropbox_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/dabio/staticdrop/dropbox"
)

func client(accessToken, url string) *dropbox.Client {
	c := &dropbox.Client{
		Config: &dropbox.Config{
			HTTPClient:  http.DefaultClient,
			AccessToken: accessToken,
			URL:         url,
		},
	}
	c.Files = &dropbox.Files{c}

	return c
}

func openFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return data
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
