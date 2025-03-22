package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIWithDefaultLogo(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Clean up any previous test output
	outputFile := "./testdata/output.html"
	os.Remove(outputFile)

	// Create configuration
	cfg := Config{
		InputDir:     "./testdata/content",
		Output:       outputFile,
		Title:        "Test Title",
		TemplatePath: "./testdata/template/test.template",
		CSSPath:      "./testdata/template/base.css",
		// LogoPath is intentionally left empty to test default behavior
	}

	err := Run(cfg)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Close the write end of the pipe to get the output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %s was not created", outputFile)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)
	t.Logf("Output file content: %s", contentStr)
	if !strings.Contains(contentStr, "Orihon") {
		t.Errorf("Output file does not contain the default logo text 'Orihon'")
	}

	// Clean up
	os.Remove(outputFile)
}

func TestCLIWithCustomLogo(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Clean up any previous test output
	outputFile := "./testdata/output-custom.html"
	os.Remove(outputFile)

	cfg := Config{
		InputDir:     "./testdata/content",
		Output:       outputFile,
		Title:        "Test Title",
		TemplatePath: "./testdata/template/test.template",
		CSSPath:      "./testdata/template/base.css",
		LogoPath:     "./testdata/custom-logo.svg",
	}

	err := Run(cfg)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Close the write end of the pipe to get the output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %s was not created", outputFile)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Custom Logo") {
		t.Errorf("Output file does not contain the custom logo text")
	}

	if !strings.Contains(contentStr, "base-css-identifier") {
		t.Errorf("Output file does not contain the base CSS identifier")
	}

	if strings.Contains(contentStr, "test-css-identifier") {
		t.Errorf("Output file contains the test CSS identifier when it should be using base.css")
	}

	// Clean up
	os.Remove(outputFile)
}

func TestCSSPathBehavior(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Clean up any previous test output
	outputFile := "./testdata/output-css.html"
	os.Remove(outputFile)

	cfg := Config{
		InputDir:     "./testdata/content",
		Output:       outputFile,
		Title:        "Test Title",
		TemplatePath: "./testdata/template/test.template",
		CSSPath:      "./testdata/template/test.css", // Using test.css instead of base.css
		LogoPath:     "./testdata/custom-logo.svg",
	}

	err := Run(cfg)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Close the write end of the pipe to get the output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %s was not created", outputFile)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "test-css-identifier") {
		t.Errorf("Output file does not contain the test CSS identifier")
	}

	if strings.Contains(contentStr, "base-css-identifier") {
		t.Errorf("Output file contains the base CSS identifier when it should be using test.css")
	}

	if !strings.Contains(contentStr, "background-color: #f5f5f5") {
		t.Errorf("Output file does not contain test.css specific styles")
	}

	// Clean up
	os.Remove(outputFile)
}

func TestLogoLinkBehavior(t *testing.T) {
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Clean up any previous test output
	outputFile := "./testdata/output-link.html"
	os.Remove(outputFile)

	testURL := "https://example.com"

	cfg := Config{
		InputDir:     "./testdata/content",
		Output:       outputFile,
		Title:        "Test Title",
		TemplatePath: "./testdata/template/test.template",
		CSSPath:      "./testdata/template/base.css",
		LogoPath:     "./testdata/custom-logo.svg",
		LogoLink:     testURL,
	}

	err := Run(cfg)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Close the write end of the pipe to get the output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file %s was not created", outputFile)
	}

	// Read the output file and check if it contains the expected logo link
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Custom Logo") {
		t.Errorf("Output file does not contain the custom logo text")
	}

	expectedLinkHTML := `<a href="https://example.com">`
	if !strings.Contains(contentStr, expectedLinkHTML) {
		t.Errorf("Output file does not contain the expected logo link. Expected: %s", expectedLinkHTML)
	}

	// Clean up
	os.Remove(outputFile)
}

func TestLogoPathDefaultBehavior(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "orihon-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create content subdirectory
	contentDir := filepath.Join(tempDir, "content")
	if err := os.Mkdir(contentDir, 0755); err != nil {
		t.Fatalf("Failed to create content directory: %v", err)
	}

	// Create a test markdown file
	mdPath := filepath.Join(contentDir, "test.md")
	if err := os.WriteFile(mdPath, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to write test markdown: %v", err)
	}

	// Create a logo.svg file in the content directory
	logoContent := `<svg><text>Content Directory Logo</text></svg>`
	logoPath := filepath.Join(contentDir, "logo.svg")
	if err := os.WriteFile(logoPath, []byte(logoContent), 0644); err != nil {
		t.Fatalf("Failed to write logo file: %v", err)
	}

	// Copy the test template to the temp directory
	templateDir := filepath.Join(tempDir, "template")
	if err := os.Mkdir(templateDir, 0755); err != nil {
		t.Fatalf("Failed to create template directory: %v", err)
	}

	// Copy test template and CSS from testdata
	templateSrc, err := os.ReadFile("./testdata/template/test.template")
	if err != nil {
		t.Fatalf("Failed to read test template: %v", err)
	}
	templateDst := filepath.Join(templateDir, "test.template")
	if err := os.WriteFile(templateDst, templateSrc, 0644); err != nil {
		t.Fatalf("Failed to write test template: %v", err)
	}

	cssSrc, err := os.ReadFile("./testdata/template/test.css")
	if err != nil {
		t.Fatalf("Failed to read test CSS: %v", err)
	}
	cssDst := filepath.Join(templateDir, "base.css")
	if err := os.WriteFile(cssDst, cssSrc, 0644); err != nil {
		t.Fatalf("Failed to write test CSS: %v", err)
	}

	outputPath := filepath.Join(tempDir, "output.html")

	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cfg := Config{
		InputDir:     contentDir,
		Output:       outputPath,
		Title:        "Test Title",
		TemplatePath: templateDst,
		CSSPath:      cssDst,
		// LogoPath is intentionally left empty to test default behavior
	}

	err = Run(cfg)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Close the write end of the pipe to get the output
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file %s was not created", outputPath)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)
	t.Logf("Output file content: %s", contentStr)
	if !strings.Contains(contentStr, "Orihon") {
		t.Errorf("Output file does not contain the default logo text 'Orihon'")
	}
}
