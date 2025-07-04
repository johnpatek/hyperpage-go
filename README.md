# hyperpage-go

![status](https://github.com/johnpatek/hyperpage-go/actions/workflows/pipeline.yml/badge.svg)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://pkg.go.dev/github.com/johnpatek/hyperpage-go)
[![codecov](https://codecov.io/gh/johnpatek/hyperpage-go/branch/master/graph/badge.svg)](https://codecov.io/gh/johnpatek/hyperpage-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/johnpatek/hyperpage-go)](https://goreportcard.com/report/github.com/johnpatek/hyperpage-go)

A pure Go implementation of Maxtek's [hyperpage](https://github.com/maxtek6/hyperpage) 
with full feature parity and interoperability.

## About

This project contains the `hyperpage` package, which can be used to store and load static 
web content from a database file. All database operations are abstracted away, with only 
their parameters exposed in the public functions.

## Usage

The following basic example demonstrates how this project could be used:

```go
package main

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/johnpatek/hyperpage-go"
)

type simplePage struct {
	path     string
	mimeType string
	content  io.Reader
}

func (p *simplePage) Path() string {
	return p.path
}

func (p *simplePage) MimeType() string {
	return p.mimeType
}

func (p *simplePage) Content() io.Reader {
	return p.content
}

func main() {
	// Open the writer
	writer, _ := hyperpage.OpenWriter(context.Background(), "hyperpage.db")
	defer writer.Close()

	// Store a page
	_ = writer.Store(context.Background(), &simplePage{
		path:     "/index.html",
		mimeType: "text/html",
		content:  strings.NewReader("<html><body><h1>Hello, World!</h1></body></html>"),
	})

	// Open the reader
	reader, _ := hyperpage.OpenReader(context.Background(), "hyperpage.db")
	defer reader.Close()

	// Load the page
	page, _ := reader.Load(context.Background(), "/index.html")
	if page != nil {
		println("Page Path:", page.Path())
		println("Page MimeType:", page.MimeType())
		println("Page Content:")
		io.Copy(os.Stdout, page.Content())
	} else {
		println("Error: page not found")
	}
}
```