# HTML2Go

A bidirectional conversion tool for translating between HTML and Go code using HTMLGo components.

## Overview

HTML2Go is a powerful utility designed to facilitate seamless conversion between:
1. Go code (using htmlgo/vuetify/vuetifyx components) to HTML
2. HTML to Go code (using htmlgo/vuetify/vuetifyx components)

This tool is particularly useful for developers working with Go-based web applications that utilize HTML component libraries like Vuetify and VuetifyX.

## Features

- **Bidirectional Conversion**: Convert in both directions - Go to HTML and HTML to Go
- **Component Support**:
  - Standard HTML elements
  - Vuetify components (`v-btn`, `v-card`, etc.)
  - VuetifyX components (`v-xbtn`, `v-xcard`, etc.)
  - Custom components
- **Attribute Handling**: Correctly maps between Go method calls and HTML attributes
- **Boolean Attributes**: Special handling for boolean attributes like `disabled`, `required`, etc.
- **Nested Structures**: Support for complex nested component structures

## Installation

### From Source

```bash
git clone https://github.com/yourusername/html2go.git
cd html2go
go build
```

### Using Go Install

```bash
go install github.com/yourusername/html2go@latest
```

## Usage

### Convert HTML to Go

```bash
# Simple usage
echo '<div>Hello, world!</div>' | html2go

# Save to file
echo '<div>Hello, world!</div>' | html2go > output.go

# Process an HTML file
cat input.html | html2go > output.go
```

### Convert Go to HTML

```bash
# Convert Go code to HTML using the -r flag
echo 'Div(Text("Hello, world!"))' | html2go -r

# Process a Go file
cat input.go | html2go -r > output.html
```

## Examples

### HTML to Go Conversion

Input (HTML):
```html
<div class="container">
  <h1>Hello, World!</h1>
  <v-btn color="primary" text="Submit"></v-btn>
</div>
```

Output (Go):
```go
package main

var n = Body(
    Div(
        Class("container"),
        H1(
            Text("Hello, World!"),
        ),
        VBtn(
            Color("primary"),
            Text("Submit"),
        ),
    ),
)
```

### Go to HTML Conversion

Input (Go):
```go
Div(
    Class("container"),
    H1(
        Text("Hello, World!"),
    ),
    VBtn(
        Color("primary"),
        Text("Submit"),
    ),
)
```

Output (HTML):
```html
<div class="container">
  <h1>Hello, World!</h1>
  <v-btn color="primary" text="Submit"></v-btn>
</div>
```

## Supported Components

- **HTML Elements**: div, span, p, h1-h6, input, button, form, etc.
- **Vuetify Components**: v-btn, v-card, v-text-field, v-select, etc.
- **VuetifyX Components**: v-xbtn, v-xcard, v-xdialog, etc.
- **Custom Components**: Support for custom-defined components

## Boolean Attributes

Boolean attributes (like `disabled`, `required`, etc.) are handled specially:
- In Go code: `Disabled(true)` 
- In HTML: `disabled` (without a value)

## Advanced Usage

### Custom Component Handling

For custom components not directly mapped, use the `Tag` function:

```go
Tag("custom-element", 
    Attr("custom-attribute", "value"),
    Text("Content"),
)
```

### Using Children Method

For components that accept child elements:

```go
Div(
    Children(
        Span(Text("Child 1")),
        Span(Text("Child 2")),
    ),
)
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
