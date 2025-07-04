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

	"github.com/johnpatek/hyperpage-go"
)

type staticPage struct {
	path     string
	mimeType string
	content  []byte
}

func (p *staticPage) Path() string {
	return p.path
}

func (p *staticPage) MimeType() string {
	return p.mimeType
}

func (p *staticPage) Content() []byte {
	return p.content
}

func main() {
	// Open the writer
	writer, _ := hyperpage.OpenWriter(context.Background(), "hyperpage.db")
	defer writer.Close()

	// Store a page
	_ = writer.Store(context.Background(), &staticPage{
		path:     "/index.html",
		mimeType: "text/html",
		content:  []byte("<html><body><h1>Hello, World!</h1></body></html>"),
	})

	// Open the reader
	reader, _ := hyperpage.OpenReader(context.Background(), "hyperpage.db")
	defer reader.Close()

	// Load the page
	page, _ := reader.Load(context.Background(), "/index.html")
	if page != nil {
		println("Page Path:", page.Path())
		println("Page MimeType:", page.MimeType())
		println("Page Content:", string(page.Content()))
	} else {
		println("Error: page not found")
	}
}
```