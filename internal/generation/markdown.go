package generation

import (
	"bytes"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var (
	mdRenderer = goldmark.New(goldmark.WithExtensions(extension.GFM))
	sanitizer  = bluemonday.UGCPolicy()
)

// Markdown renders CommonMark/GFM markdown to sanitized HTML.
func Markdown(md string) (string, error) {
	var buf bytes.Buffer
	if err := mdRenderer.Convert([]byte(md), &buf); err != nil {
		return "", err
	}
	return string(sanitizer.SanitizeBytes(buf.Bytes())), nil
}
