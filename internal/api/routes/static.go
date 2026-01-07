//go:build !debug
// +build !debug

package routes

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/akramboussanni/marchive/config"
	"github.com/akramboussanni/marchive/internal/model"
	"github.com/akramboussanni/marchive/internal/repo"
	"github.com/go-chi/chi/v5"
)

var (
	globalRepos     *repo.Repos
	cachedIndexHTML string
	indexHTMLMutex  sync.RWMutex
	bookHashRegex   = regexp.MustCompile(`^/(book|read)/([a-fA-F0-9]+)$`)
)

func setupStaticRoutes(r chi.Router, repos *repo.Repos) {
	globalRepos = repos

	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		frontendDir = "./frontend/build"
	}

	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		return
	}

	// Load and cache index.html once at startup
	indexPath := filepath.Join(frontendDir, "index.html")
	if err := loadIndexHTML(indexPath); err != nil {
		// Log error but continue - will try to read on each request as fallback
		fmt.Printf("Warning: failed to cache index.html: %v\n", err)
	}

	// Create a file server for the frontend directory
	fs := http.FileServer(http.Dir(frontendDir))

	// Serve static files (assets, favicon, etc.)
	r.Get("/assets/*", func(w http.ResponseWriter, req *http.Request) {
		fs.ServeHTTP(w, req)
	})

	// Serve favicon and other root-level static files
	r.Get("/favicon.ico", func(w http.ResponseWriter, req *http.Request) {
		fs.ServeHTTP(w, req)
	})

	// Catch-all: serve index.html for SPA routing
	r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		// Skip API routes
		if strings.HasPrefix(req.URL.Path, "/api/") {
			http.NotFound(w, req)
			return
		}

		// Check if the requested file exists
		requestedPath := filepath.Join(frontendDir, req.URL.Path)
		if info, err := os.Stat(requestedPath); err == nil && !info.IsDir() {
			// File exists, serve it directly
			fs.ServeHTTP(w, req)
			return
		}

		// Otherwise serve index.html for SPA routing with dynamic meta tags
		serveIndexWithMeta(w, req)
	})
}

// loadIndexHTML loads and caches the index.html file
func loadIndexHTML(indexPath string) error {
	htmlBytes, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	indexHTMLMutex.Lock()
	cachedIndexHTML = string(htmlBytes)
	indexHTMLMutex.Unlock()

	return nil
}

// serveIndexWithMeta serves index.html with dynamic OpenGraph meta tags for book pages
func serveIndexWithMeta(w http.ResponseWriter, r *http.Request) {
	// Get cached HTML
	indexHTMLMutex.RLock()
	html := cachedIndexHTML
	indexHTMLMutex.RUnlock()

	// Fallback: if cache is empty, try to read file
	if html == "" {
		frontendDir := os.Getenv("FRONTEND_DIR")
		if frontendDir == "" {
			frontendDir = "./frontend/build"
		}
		indexPath := filepath.Join(frontendDir, "index.html")
		htmlBytes, err := os.ReadFile(indexPath)
		if err != nil {
			http.Error(w, "index.html not found", http.StatusNotFound)
			return
		}
		html = string(htmlBytes)
	}

	// Check if this is a book detail or read page
	matches := bookHashRegex.FindStringSubmatch(r.URL.Path)

	if len(matches) == 3 {
		pageType := matches[1] // "book" or "read"
		hash := matches[2]

		// Get book data from database
		if globalRepos != nil && globalRepos.Book != nil {
			book, err := globalRepos.Book.GetBookByHash(context.Background(), hash)
			if err == nil && book != nil {
				// Generate full URL for the page
				scheme := "https"
				if !config.App.TLSEnabled {
					scheme = "http"
				}
				pageURL := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path)

				// Inject OpenGraph meta tags
				metaTags := generateBookMetaTags(book, pageType, r.Host, pageURL)
				html = strings.Replace(html, "<!--DYNAMIC_META-->", metaTags, 1)
			}
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

// generateBookMetaTags creates OpenGraph meta tags for a book
func generateBookMetaTags(book *model.SavedBook, pageType, host, pageURL string) string {
	if book == nil {
		return ""
	}

	authorInfo := ""
	if book.Authors != "" {
		authorInfo = " by " + book.Authors
	}

	var pageTitle, ogTitle, ogDescription string
	if pageType == "read" {
		pageTitle = fmt.Sprintf("Read %s%s - mArchive", book.Title, authorInfo)
		ogTitle = "Read " + book.Title + " on mArchive"
		ogDescription = book.Title + authorInfo + " - Read online on mArchive."
	} else {
		pageTitle = fmt.Sprintf("%s%s - mArchive", book.Title, authorInfo)
		ogTitle = book.Title + authorInfo + " | mArchive"
		ogDescription = book.Title + authorInfo + " - " + book.Format + " - " + book.Size + "."
	}

	// Determine cover image URL - must be absolute public URL for scrapers
	var coverImageURL string
	if book.CoverURL != "" && (strings.HasPrefix(book.CoverURL, "http://") || strings.HasPrefix(book.CoverURL, "https://")) {
		// Already an absolute URL
		coverImageURL = book.CoverURL
	} else if book.CoverData != "" && (strings.HasPrefix(book.CoverData, "http://") || strings.HasPrefix(book.CoverData, "https://")) {
		// CoverData is an absolute URL
		coverImageURL = book.CoverData
	}
	// Note: We skip data: URLs because most scrapers don't support them

	var buf bytes.Buffer

	// Page title
	buf.WriteString("<title>")
	buf.WriteString(template.HTMLEscapeString(pageTitle))
	buf.WriteString("</title>\n    ")

	// Meta description
	buf.WriteString(`<meta name="description" content="`)
	buf.WriteString(template.HTMLEscapeString(ogDescription))
	buf.WriteString(`">`)
	buf.WriteString("\n    ")

	// OpenGraph tags
	buf.WriteString(`<meta property="og:title" content="`)
	buf.WriteString(template.HTMLEscapeString(ogTitle))
	buf.WriteString(`">`)
	buf.WriteString("\n    ")

	buf.WriteString(`<meta property="og:description" content="`)
	buf.WriteString(template.HTMLEscapeString(ogDescription))
	buf.WriteString(`">`)
	buf.WriteString("\n    ")

	buf.WriteString(`<meta property="og:url" content="`)
	buf.WriteString(template.HTMLEscapeString(pageURL))
	buf.WriteString(`">`)
	buf.WriteString("\n    ")

	// Set og:type to "book" for book pages
	buf.WriteString(`<meta property="og:type" content="book">`)
	buf.WriteString("\n    ")

	if coverImageURL != "" {
		buf.WriteString(`<meta property="og:image" content="`)
		buf.WriteString(template.HTMLEscapeString(coverImageURL))
		buf.WriteString(`">`)
		buf.WriteString("\n    ")

		// Twitter card with large image
		buf.WriteString(`<meta name="twitter:card" content="summary_large_image">`)
		buf.WriteString("\n    ")

		buf.WriteString(`<meta name="twitter:image" content="`)
		buf.WriteString(template.HTMLEscapeString(coverImageURL))
		buf.WriteString(`">`)
		buf.WriteString("\n    ")
	}

	// Twitter meta tags
	buf.WriteString(`<meta name="twitter:title" content="`)
	buf.WriteString(template.HTMLEscapeString(ogTitle))
	buf.WriteString(`">`)
	buf.WriteString("\n    ")

	buf.WriteString(`<meta name="twitter:description" content="`)
	buf.WriteString(template.HTMLEscapeString(ogDescription))
	buf.WriteString(`">`)

	return buf.String()
}
