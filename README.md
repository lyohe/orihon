# Orihon

Orihon is a simple static site generator written in Go, designed to compile Markdown files into a single-page HTML document with a built-in Table of Contents.

![Screenshot](https://github.com/user-attachments/assets/5109c8bd-cc18-419d-9965-20bf45a0ee25)

<p align="center">
  <img src="https://github.com/user-attachments/assets/f606da9b-1ece-4761-bcf1-304db61b8f25" width="48%" />
  <img src="https://github.com/user-attachments/assets/174beb5a-2144-4a83-b3de-360a4c8f18af" width="48%" />
</p>

## Features

- **Single-Page HTML Generation**  
  Compiles Markdown files into a clean, single HTML page.

- **Automatic Table of Contents**  
  Automatically extracts headings from Markdown files to generate a navigable sidebar.

- **Responsive and User-Friendly**  
  Built-in support for responsive design, dark mode, and interactive sidebar navigation.

## Usage

### Prerequisites

To build and run Orihon, ensure your environment meets the following requirements:

- **Operating System:** Linux, macOS, or Windows
- **Go:** Version 1.24 or later ([Installation Instructions](https://go.dev/doc/install))
- **Make:** GNU Make installed on your system (for Windows, consider using [Chocolatey](https://chocolatey.org/install) or [WSL](https://learn.microsoft.com/windows/wsl/install))

### Generate HTML

1. Place your Markdown files in the `content/` directory.
2. _(Optional)_ Adjust the HTML template and CSS in the `template/` directory.
3. Run the following command:

```bash
make
```

This command generates an `index.html` file in the project root.

### Custom Generation Options

You can manually invoke the generator with additional parameters:

```sh
go run ./cmd/main.go \
  -input=./content \
  -output=./index.html \
  -title="My Documentation" \
  -template=./template/base.template \
  -css=./template/base.css \
  -logo=./path/to/custom-logo.svg \
  -logo-link="https://example.com"
```

| Option       | Default                                 | Description                                                                          |
| ------------ | --------------------------------------- | ------------------------------------------------------------------------------------ |
| `-input`     | `./content`                             | Path to Markdown input directory                                                     |
| `-output`    | `index.html`                            | Output HTML file                                                                     |
| `-title`     | `"Orihon Docs - Single Page Generator"` | Title for the generated HTML                                                         |
| `-template`  | `./template/base.template`              | Path to custom HTML template                                                         |
| `-css`       | `./template/base.css`                   | Path to CSS stylesheet                                                               |
| `-logo`      | `"./content/logo.svg"`                  | Logo SVG file path (if not specified, defaults to 'logo.svg' in the input directory) |
| `-logo-link` | `"https://github.com/lyohe/orihon"`     | URL to navigate to when the logo is clicked                                          |

### Customize Templates and Styles

You can edit base.template and base.css to customize HTML generated.

## License

MIT Licensed
