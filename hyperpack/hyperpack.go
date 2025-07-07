package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnpatek/hyperpage-go"
	"github.com/spf13/pflag"
)

type mappedPage struct {
	path     string
	mimeType string
	content  io.Reader
}

func (p *mappedPage) Path() string {
	return p.path
}

func (p *mappedPage) MimeType() string {
	return p.mimeType
}

func (p *mappedPage) Content() io.Reader {
	return p.content
}

func newMappedPage(base, path string) (hyperpage.Page, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	webPath := strings.ReplaceAll(strings.TrimPrefix(path, base), string(os.PathSeparator), "/")
	if !strings.HasPrefix(webPath, "/") {
		webPath = "/" + webPath
	}
	mime := mime.TypeByExtension(filepath.Ext(webPath))

	return &mappedPage{
		path:     webPath,
		mimeType: mime,
		content:  bytes.NewReader(content),
	}, nil
}

func main() {
	var (
		output  string
		verbose bool
	)

	pflag.StringVarP(&output, "output", "o", "hyperpage.db", "Output database file")
	pflag.BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	pflag.Parse()

	if pflag.NArg() < 1 {
		fmt.Println("Missing required directory argument.")
		os.Exit(1)
	}
	dir := pflag.Arg(0)

	writer, err := hyperpage.OpenWriter(context.Background(), output)
	if err != nil {
		panic(err)
	}
	defer writer.Close()

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}
		page, err := newMappedPage(dir, path)
		if err != nil {
			return err
		}
		if verbose {
			fmt.Printf("Adding: %s [%s]\n", page.Path(), page.MimeType())
		}
		return writer.Store(context.Background(), page)
	})
	if err != nil {
		panic(err)
	}
}
