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
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	return db, mock
}

func TestSqlOpen(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(0))
	defer cancel()
	db, err := sqlOpen(ctx, ":memory:", "")
	assert.Error(t, err)
	assert.Nil(t, db)
	db, err = sqlOpen(context.Background(), ":memory:", "some invalid query")
	assert.Error(t, err)
	assert.Nil(t, db)
	db, err = sqlOpen(context.Background(), ":memory:", "")
	assert.NoError(t, err)
	assert.NotNil(t, db)
	defer db.Close()
}

func TestOpenReader(t *testing.T) {
	// Test opening a reader with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(0))
	defer cancel()
	reader, err := OpenReader(ctx, ":memory:")
	assert.Error(t, err)
	assert.Nil(t, reader)

	// Test opening a reader successfully
	reader, err = OpenReader(context.Background(), ":memory:")
	assert.NoError(t, err)
	assert.NotNil(t, reader)
	defer reader.Close()
}

func TestOpenWriter(t *testing.T) {
	// Test opening a writer with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(0))
	defer cancel()
	writer, err := OpenWriter(ctx, ":memory:")
	assert.Error(t, err)
	assert.Nil(t, writer)

	// Test opening a writer successfully
	writer, err = OpenWriter(context.Background(), ":memory:")
	assert.NoError(t, err)
	assert.NotNil(t, writer)
	defer writer.Close()
}

func TestStorePage(t *testing.T) {
	db, mock := newMock()
	defer db.Close()
	writer := &Writer{db: db}

	// Test storing a nil page
	err := writer.Store(context.Background(), nil)
	assert.Error(t, err)

	// Test timeout error
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(0))
	defer cancel()
	err = writer.Store(ctx, &sqlitePage{
		path:     "/test/path",
		mimeType: "text/plain",
		content:  []byte("test content"),
	})
	assert.Error(t, err)
	// Test successful storage
	mock.ExpectExec(regexp.QuoteMeta(insertPageQuery)).WithArgs("/test/path", "text/plain", []byte("test content")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err = writer.Store(context.Background(), &sqlitePage{
		path:     "/test/path",
		mimeType: "text/plain",
		content:  []byte("test content"),
	})
	assert.NoError(t, err)
}

func TestLoadPage(t *testing.T) {
	db, mock := newMock()
	defer db.Close()
	reader := &Reader{db: db}

	// Test loading a page with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(0))
	defer cancel()
	page, err := reader.Load(ctx, "/test/path")
	assert.Error(t, err)
	assert.Nil(t, page)

	// Test loading a page successfully
	mock.ExpectQuery(regexp.QuoteMeta(selectPageQuery)).WithArgs("/test/path").
		WillReturnRows(sqlmock.NewRows([]string{"mime_type", "content"}).
			AddRow("text/plain", []byte("test content")))
	page, err = reader.Load(context.Background(), "/test/path")
	assert.NoError(t, err)
	assert.NotNil(t, page)
	assert.Equal(t, "/test/path", page.Path())
	assert.Equal(t, "text/plain", page.MimeType())
	assert.Equal(t, []byte("test content"), page.Content())
}
