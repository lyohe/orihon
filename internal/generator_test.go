package orihon

import (
	"strings"
	"testing"
)

// TestGenerateHTML verifies the core functionality of the HTML generation process.
// It tests that the generator correctly:
// - Processes Markdown into HTML
// - Creates a proper table of contents structure
// - Applies the HTML template with the correct title
func TestGenerateHTML(t *testing.T) {
	md := `
# Title H1
## Subtitle H2
Some paragraph text.

### Sub-subtitle H3
Another paragraph.

`
	css := `body { font-family: sans-serif; }`
	tpl := `
<html>
<head><title>{{.Title}}</title><style>{{.InlineCss}}</style></head>
<body>
<div class="sidebar">{{.Sidebar}}</div>
<div class="content">{{.Content}}</div>
</body>
</html>
`

	cfg := Config{
		Title: "Test Orihon",
	}

	html, err := GenerateHTML([]byte(md), []byte(css), []byte(tpl), cfg)
	if err != nil {
		t.Fatalf("GenerateHTML error: %v", err)
	}

	if !strings.Contains(html, "<title>Test Orihon</title>") {
		t.Error("output HTML does not contain correct <title>")
	}
	if !strings.Contains(html, "Title H1") {
		t.Error("output HTML does not contain heading text")
	}
	if !strings.Contains(html, "<ul class=\"toc\">") {
		t.Error("output HTML does not contain TOC <ul> structure")
	}
}
