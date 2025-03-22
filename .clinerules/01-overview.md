# Orihon Development Guide

## Project Overview

Orihon is a simple static site generator written in Go, designed to compile a Markdown file into a single-page HTML document with a built-in Table of Contents.

## Goals

The goal of this project is to develop a static site generator that combines multiple Markdown files into a single-page HTML document.

## Architecture

The project will be developed using the following technologies:

- Go programming language
- HTML/CSS/JavaScript
- Markdown

A Makefile is used to automate the build process.

The Go program reads Markdown files, a template file, and CSS files from specified directories, parses them, and then generates a single-page HTML file.

The generated HTML will include styles from the CSS files and use the templates for layout.

## Definitions

- Core Generator: All logic within the internal/ directory responsible for parsing Markdown and generating HTML.
- Templates and Assets: All files located in the template/ directory, such as HTML templates and CSS stylesheets.
- Content Directory: The content/ directory contains Markdown input files and the logo SVG file.

## Coding Rules

### File Organization

- All Go logic for HTML generation must be placed within the internal/ directory.
- CLI logic and argument parsing should reside in cmd/orihon/.
- Templates, CSS, and assets (such as SVG logos) should be stored appropriately within the template/ or content/ directories.

### Directory Structure

The project follows this directory structure:

```
.
├── cmd/                        # Command-line interface
│   ├── main.go                 # CLI entry point
│   └── testdata/               # Test fixtures
├── content/                    # Markdown content files
├── internal/                   # Core generator logic
└── template/                   # Templates and assets
```

### Template Customization

- All HTML markup must be contained exclusively within template/base.template.
- CSS styles should only be modified in template/base.css.
- SVG logos should be placed in content/logo.svg.

## HTML Generation

The generator logic (internal/generator.go) should follow these steps:

1. Parse Markdown content using Goldmark.
2. Extract headings and build a structured Table of Contents (TOC).
3. Dynamically load the logo SVG file (if available).
4. Generate the final HTML by combining the parsed Markdown content, CSS styles, SVG logo, and provided template.

## Testing

- Write tests in the same directory as implementation files (e.g., internal/generator_test.go).
- Ensure test cases cover Markdown parsing, heading extraction, and template rendering.

## Error Handling

- Gracefully handle file I/O errors with clear log messages.
- Use log.Fatalf for unrecoverable errors.
