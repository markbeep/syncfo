package generation

import (
	"strings"
	"testing"
)

func TestMarkdownRendersGFM(t *testing.T) {
	html, err := Markdown("~~gone~~\n\n```go\nfmt.Println(1)\n```\n")
	if err != nil {
		t.Fatalf("Markdown error: %v", err)
	}
	if !strings.Contains(html, "<del>") {
		t.Errorf("strikethrough not rendered: %s", html)
	}
	if !strings.Contains(html, "<pre><code") {
		t.Errorf("code block not rendered: %s", html)
	}
}

func TestMarkdownSanitizesScriptAndURLs(t *testing.T) {
	input := "<script>alert(1)</script>\n\n<img src=x onerror=alert(1)>\n\n[click](javascript:alert(2))\n"
	html, err := Markdown(input)
	if err != nil {
		t.Fatalf("Markdown error: %v", err)
	}
	for _, banned := range []string{"<script", "<img", "javascript:"} {
		if strings.Contains(html, banned) {
			t.Errorf("output contains %q: %s", banned, html)
		}
	}
}

func TestMarkdownStripsEventHandlers(t *testing.T) {
	got := string(sanitizer.SanitizeBytes([]byte(`<a href="/" onerror="alert(1)" onclick="alert(2)">x</a>`)))
	for _, banned := range []string{"onerror", "onclick"} {
		if strings.Contains(got, banned) {
			t.Errorf("output contains %q: %s", banned, got)
		}
	}
}
