package page

import (
	"fmt"
	"strings"
	"time"
)

// URL joins a URL base path and a relative path with exactly one slash
// between them. A leading slash is always produced, so results are absolute
// within the site.
func URL(base, path string) string {
	base = strings.TrimRight(base, "/")
	path = strings.TrimLeft(path, "/")
	if base == "" {
		return "/" + path
	}
	return base + "/" + path
}

// ListURL returns the URL for a list page. Page 0 is the temporary newest
// page; stable pages are 1-based.
func ListURL(base string, number int) string {
	return URL(base, fmt.Sprintf("list/%d/", number))
}

// ThreadURL returns the URL for page 1 of a thread (no page suffix).
func ThreadURL(base string, id int64) string {
	return URL(base, fmt.Sprintf("thread/%d", id))
}

// ThreadPageURL returns the URL for a message page of a thread. Page 1 is the
// bare thread URL; later pages carry a numeric suffix.
func ThreadPageURL(base string, id int64, page int) string {
	if page <= 1 {
		return ThreadURL(base, id)
	}
	return URL(base, fmt.Sprintf("thread/%d/%d", id, page))
}

// FormatTime renders a timestamp in a compact UTC form for display.
func FormatTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04")
}
