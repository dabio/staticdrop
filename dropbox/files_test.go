package dropbox_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dabio/staticdrop/dropbox"
)

func TestListFolderSuccess(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(openFile("testdata/list_folder_success.json"))
		}),
	)
	defer server.Close()
	client := client("token", server.URL)

	in := &dropbox.ListFolderArg{}
	out, err := client.Files.ListFolder(in)
	ok(t, err)

	t.Run("list folder result", func(t *testing.T) {
		equals(t, false, out.HasMore)
		equals(t, "fcce3e7dc8d3", out.Cursor)
		equals(t, 2, len(out.Entries))
	})

	t.Run("list folder first entry (folder)", func(t *testing.T) {
		r := out.Entries[0]

		equals(t, "folder", r.Tag)
		equals(t, "staticdrop", r.Name)
		equals(t, "/staticdrop", r.PathLower)
		equals(t, "/staticdrop", r.PathDisplay)
		equals(t, "id:f1f9321c-2c9c-4b86-904b-1b11db303a18", r.ID)

		equals(t, true, r.ClientModified.IsZero())
		equals(t, true, r.ServerModified.IsZero())
		equals(t, "", r.Rev)
		equals(t, uint64(0), r.Size)
		equals(t, "", r.ContentHash)
	})

	t.Run("list folder second entry (file)", func(t *testing.T) {
		r := out.Entries[1]

		equals(t, "file", r.Tag)
		equals(t, "moving_cli.png", r.Name)
		equals(t, "/staticdrop/moving_cli.png", r.PathLower)
		equals(t, "/staticdrop/moving_cli.png", r.PathDisplay)
		equals(t, "id:878fc37b-2b70-488c-ab02-fb5d01b179ec", r.ID)

		// 2018-02-07T08:45:34Z
		equals(t, time.Date(2018, 2, 7, 8, 45, 34, 0, time.UTC), r.ClientModified)
		// 2018-02-07T08:45:35Z
		equals(t, time.Date(2018, 2, 7, 8, 45, 35, 0, time.UTC), r.ServerModified)
		equals(t, "28c2ab9d0", r.Rev)
		equals(t, uint64(50647), r.Size)
		equals(t, "424da3cf0c0b67fb71ab74ab7a8e4ad622a5d363a6637af79b64dd769fccecd8", r.ContentHash)
	})
}
