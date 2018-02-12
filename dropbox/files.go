package dropbox

import (
	"encoding/json"
	"time"
)

const listFolderURL = "/files/list_folder"

// Files client for files and folders.
type Files struct {
	*Client
}

// NewFiles client.
func NewFiles(config *Config) *Files {
	return &Files{
		Client: &Client{
			Config: config,
		},
	}
}

// Metadata for files or folders.
type Metadata struct {
	// folder, file:
	Tag         string `json:".tag"`
	Name        string `json:"name"`
	PathLower   string `json:"path_lower"`
	PathDisplay string `json:"path_display"`
	ID          string `json:"id"`

	// file:
	ClientModified time.Time `json:"client_modified"`
	ServerModified time.Time `json:"server_modified"`
	Rev            string    `json:"rev"`
	Size           uint64    `json:"size"`
	ContentHash    string    `json:"content_hash"`
}

// ListFolderArg are the valid parameters for the API.
type ListFolderArg struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive"`
}

// ListFolderResult is the result returned by the API.
type ListFolderResult struct {
	Entries []*Metadata `json:"entries"`
	Cursor  string      `json:"cursor"`
	HasMore bool        `json:"has_more"`
}

// ListFolder starts returning the contents of a folder.
// https://www.dropbox.com/developers/documentation/http/documentation#files-list_folder
func (c *Files) ListFolder(in *ListFolderArg) (out *ListFolderResult, err error) {
	body, err := c.call(listFolderURL, in)
	if err != nil {
		return
	}
	defer body.Close()

	err = json.NewDecoder(body).Decode(&out)
	return
}
