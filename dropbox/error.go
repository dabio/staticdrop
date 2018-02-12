package dropbox

// Error response of API.
type Error struct {
	Status     string
	StatusCode int
	Summary    string `json:"error_summary"`
}

func (e *Error) Error() string {
	return e.Summary
}
