package orihon

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type Heading struct {
	Level    int
	Text     string
	AnchorID string
	Children []*Heading
}

type Config struct {
	Title     string
	Sidebar   template.HTML
	Content   template.HTML
	InlineCss template.CSS
	LogoHTML  template.HTML
	LogoLink  string
	LogoPath  string // Path to the logo SVG file
}

func GenerateHTML(
	mdContent []byte,
	cssContent []byte,
	tmplContent []byte,
	cfg Config,
) (string, error) {
	logoHTML, err := loadLogo(cfg.LogoPath)
	if err != nil {
		return "", err
	}
	cfg.LogoHTML = logoHTML

	// Convert Markdown content to HTML
	gm := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
	doc := gm.Parser().Parse(text.NewReader(mdContent))

	headings, err := extractHeadings(doc, mdContent)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := gm.Convert(mdContent, &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown: %w", err)
	}
	htmlContent := buf.String()

	sidebarHTML := buildSidebarHTML(headings)

	tpl, err := template.New("base").Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	cfg.Sidebar = template.HTML(sidebarHTML)
	cfg.Content = template.HTML(htmlContent)
	cfg.InlineCss = template.CSS(cssContent)

	var outBuf bytes.Buffer
	if err := tpl.Execute(&outBuf, cfg); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return outBuf.String(), nil
}

// loadLogo attempts to load the logo SVG file from the specified path.
func loadLogo(path string) (template.HTML, error) {
	logoSVG, err := os.ReadFile(path)
	if err == nil {
		return template.HTML(string(logoSVG)), nil
	}

	defaultLogo := template.HTML(`
		<svg xmlns="http://www.w3.org/2000/svg" width="140" height="28" viewBox="0 0 140 28">
			<path d="M2 24 L7 4 L12 24 L17 4 L22 24 L27 4" fill="none" stroke="#4A90E2" stroke-width="2" stroke-linejoin="round"/>
			<text x="36" y="20" font-family="Arial, sans-serif" font-size="16" fill="#333">Orihon</text>
		</svg>`)

	return defaultLogo, nil
}

// extractHeadings parses headings (H1..H6) from Markdown AST and returns them in a hierarchical structure
func extractHeadings(doc ast.Node, source []byte) ([]*Heading, error) {
	var flat []*Heading

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		h, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}

		var text string
		for n := h.FirstChild(); n != nil; n = n.NextSibling() {
			if t, ok := n.(*ast.Text); ok {
				text += string(t.Segment.Value(source))
			}
		}
		anchor, ok := h.AttributeString("id")
		if !ok {
			return ast.WalkStop, fmt.Errorf("heading(%s) does not have id attribute", text)
		}

		flat = append(flat, &Heading{
			Level:    h.Level,
			Text:     text,
			AnchorID: string(anchor.([]byte)),
		})

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, err
	}

	return buildHeadingTree(flat), nil
}

// buildHeadingTree organizes flat headings into a nested tree structure based on their Level
func buildHeadingTree(flat []*Heading) []*Heading {
	if len(flat) == 0 {
		return nil
	}
	var root []*Heading
	stack := []*Heading{}

	for _, h := range flat {
		if len(stack) == 0 {
			stack = append(stack, h)
			root = append(root, h)
			continue
		}
		top := stack[len(stack)-1]

		if h.Level > top.Level {
			top.Children = append(top.Children, h)
			stack = append(stack, h)
		} else {
			for len(stack) > 0 && stack[len(stack)-1].Level >= h.Level {
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				root = append(root, h)
				stack = append(stack, h)
			} else {
				newTop := stack[len(stack)-1]
				newTop.Children = append(newTop.Children, h)
				stack = append(stack, h)
			}
		}
	}

	return root
}

// buildSidebarHTML generates an HTML string with nested <ul> elements from the heading tree
func buildSidebarHTML(headings []*Heading) string {
	var sb strings.Builder
	sb.WriteString("<ul class=\"toc\">\n")
	for _, h := range headings {
		buildSidebarItem(&sb, h, 0)
	}
	sb.WriteString("</ul>\n")
	return sb.String()
}

// buildSidebarItem recursively constructs the HTML for a single sidebar item and its children
func buildSidebarItem(sb *strings.Builder, h *Heading, depth int) {
	indentClass := fmt.Sprintf("indent-%d", depth)
	hasChildren := len(h.Children) > 0

	sb.WriteString(fmt.Sprintf("<li class=\"toc-item %s\">", indentClass))

	if hasChildren {
		sb.WriteString("<details>\n")
		sb.WriteString(fmt.Sprintf("<summary><a href=\"#%s\">%s</a></summary>\n", h.AnchorID, h.Text))
		sb.WriteString("<ul>\n")
		for _, child := range h.Children {
			buildSidebarItem(sb, child, depth+1)
		}
		sb.WriteString("</ul>\n</details>\n")
	} else {
		sb.WriteString(fmt.Sprintf("<a href=\"#%s\">%s</a>", h.AnchorID, h.Text))
	}

	sb.WriteString("</li>\n")
}
