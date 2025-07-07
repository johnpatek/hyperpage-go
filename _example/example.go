package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/maxtek6/sigfn-go"

	"github.com/johnpatek/hyperpage-go"
)

type pageHandler struct {
	reader *hyperpage.Reader
}

func newPageHandler(reader *hyperpage.Reader) *pageHandler {
	dbPath := path.Join(getDirectory(), "hyperpage.db")
	reader, err := hyperpage.OpenReader(context.Background(), dbPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to open hyperpage reader: %v", err))
	}
	return &pageHandler{
		reader: reader,
	}
}

func (h *pageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}
	h.servePage(w, path)
}

func (h *pageHandler) servePage(w http.ResponseWriter, path string) {
	page, err := h.reader.Load(context.Background(), path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading page: %v", err), http.StatusInternalServerError)
		return
	}
	if page == nil {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", page.MimeType())
	_, _ = io.Copy(w, page.Content())
}

func getDirectory() string {
	dir, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("Failed to get current directory: %v", err))
	}
	return filepath.Dir(dir)
}

func main() {
	handler := newPageHandler(nil) // Replace nil with actual reader initialization
	server := &http.Server{
		Addr:    ":12345",
		Handler: handler,
	}

	signalHandler := func(sig os.Signal) {
		fmt.Println("Received signal:", sig)
		_ = server.Shutdown(context.Background())
	}

	sigfn.Handle(syscall.SIGINT, signalHandler)
	_ = server.ListenAndServe()
	fmt.Println("Server stopped gracefully")
}
