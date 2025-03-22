// filename: cmd/orihon/main.go
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	orihon "github.com/lyohe/orihon/internal"
)

type Config struct {
	InputDir     string
	Output       string
	Title        string
	TemplatePath string
	CSSPath      string
	LogoPath     string
	LogoLink     string
}

func ParseFlags() Config {
	var cfg Config

	flag.StringVar(&cfg.InputDir, "input", "./content", "Directory containing Markdown files")
	flag.StringVar(&cfg.Output, "output", "index.html", "Output HTML filename")
	flag.StringVar(&cfg.Title, "title", "Orihon Docs - Single Page Generator", "HTML title")
	flag.StringVar(&cfg.TemplatePath, "template", "./template/base.template", "HTML template file path")
	flag.StringVar(&cfg.CSSPath, "css", "./template/base.css", "CSS file path")
	flag.StringVar(&cfg.LogoPath, "logo", "", "Logo SVG file path (defaults to 'logo.svg' in input directory if exists)")
	flag.StringVar(&cfg.LogoLink, "logo-link", "https://github.com/lyohe/orihon", "URL to navigate to when the logo is clicked")
	flag.Parse()

	if cfg.LogoPath == "" {
		cfg.LogoPath = filepath.Join(cfg.InputDir, "logo.svg")
	}

	return cfg
}

func Run(cfg Config) error {
	// Recursively collect all Markdown content from the input directory
	var mdBuffer bytes.Buffer
	err := filepath.Walk(cfg.InputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".md" {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read Markdown file %s: %w", path, err)
			}
			mdBuffer.Write(content)
			mdBuffer.WriteString("\n\n")
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading Markdown files: %w", err)
	}
	mdBytes := mdBuffer.Bytes()

	// Load template files
	cssBytes, err := os.ReadFile(cfg.CSSPath)
	if err != nil {
		return fmt.Errorf("failed to read CSS file: %w", err)
	}
	templateBytes, err := os.ReadFile(cfg.TemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Process Markdown content to build the output HTML structure
	orihonCfg := orihon.Config{
		Title:    cfg.Title,
		LogoPath: cfg.LogoPath,
		LogoLink: cfg.LogoLink,
	}
	outHTML, err := orihon.GenerateHTML(mdBytes, cssBytes, templateBytes, orihonCfg)
	if err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	// Write the final HTML to the output file
	outFile, err := os.Create(cfg.Output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	if _, err := outFile.Write([]byte(outHTML)); err != nil {
		return fmt.Errorf("failed to write HTML to file: %w", err)
	}

	fmt.Printf("Done! Generated -> %s\n", cfg.Output)
	return nil
}

func main() {
	cfg := ParseFlags()
	if err := Run(cfg); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
