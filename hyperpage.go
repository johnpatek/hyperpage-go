// Copyright (c) 2024 John R Patek Sr
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package hyperpage

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var (
	createTableQuery string = `CREATE TABLE IF NOT EXISTS hyperpage (
        path TEXT PRIMARY KEY,
        mime_type TEXT,
        content BLOB);
        CREATE INDEX IF NOT EXISTS path_index ON hyperpage (path);`
	insertPageQuery string = `INSERT OR REPLACE INTO hyperpage (path, mime_type, content) VALUES (?, ?, ?);`
	selectPageQuery string = `SELECT mime_type, content FROM hyperpage WHERE path = ?;`
)

// Page represents a single page in the hyperpage database
type Page interface {
	Path() string
	MimeType() string
	Content() []byte
}

type sqlitePage struct {
	path     string
	mimeType string
	content  []byte
}

// Reader is used to load pages from the hyperpage database
type Reader struct {
	db *sql.DB
}

// Writer is used to store pages in the hyperpage database
type Writer struct {
	db *sql.DB
}

func (p *sqlitePage) Path() string {
	return p.path
}

func (p *sqlitePage) MimeType() string {
	return p.mimeType
}

func (p *sqlitePage) Content() []byte {
	return p.content
}

// OpenReader opens a reader for the hyperpage database at the specified path
func OpenReader(ctx context.Context, path string) (*Reader, error) {
	db, err := sqlOpen(ctx, path, "")
	if err != nil {
		return nil, fmt.Errorf("hyperpage.OpenReader: %v", err)
	}
	return &Reader{
		db: db,
	}, nil
}

// Close closes the database connection for the reader
func (r *Reader) Close() {
	_ = r.db.Close()
}

// Load retrieves a page from the hyperpage database by its path
func (r *Reader) Load(ctx context.Context, path string) (Page, error) {
	row := r.db.QueryRowContext(ctx, selectPageQuery, path)
	var mimeType string
	var content []byte
	err := row.Scan(&mimeType, &content)
	if err != nil {
		return nil, fmt.Errorf("hyperpage.Load: %v", err)
	}
	return &sqlitePage{
		path:     path,
		mimeType: mimeType,
		content:  content,
	}, nil
}

// OpenWriter opens a writer for the hyperpage database at the specified path
func OpenWriter(ctx context.Context, path string) (*Writer, error) {
	db, err := sqlOpen(ctx, path, createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("hyperpage.OpenWriter: %v", err)
	}
	return &Writer{
		db: db,
	}, nil
}

// Close closes the database connection for the writer
func (w *Writer) Close() {
	_ = w.db.Close()
}

// Store saves a page to the hyperpage database
func (w *Writer) Store(ctx context.Context, page Page) error {
	if page == nil {
		return fmt.Errorf("hyperpage.Store: cannot store nil page")
	}
	_, err := w.db.ExecContext(ctx, insertPageQuery, page.Path(), page.MimeType(), page.Content())
	if err != nil {
		return fmt.Errorf("hyperpage.Store: %v", err)
	}
	return nil
}

func sqlOpen(ctx context.Context, path, firstExec string) (*sql.DB, error) {
	db, _ := sql.Open("sqlite", path)
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	if firstExec != "" {
		_, err = db.ExecContext(ctx, firstExec)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
