package server

import (
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"
)

// handleStatic serves the pre-rendered SSG output.
//
// Resolution order for a request path P:
//  1. If P has a file extension, serve that exact file (assets).
//  2. Otherwise serve P/index.html (vite-react-ssg emits one folder per route).
//  3. On miss, serve dist/404.html with a 404 status.
//
// Fingerprinted assets under /assets/ get an immutable, one-year cache; HTML
// is always revalidated so content updates go live immediately.
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clean := path.Clean(r.URL.Path)
	hasExt := path.Ext(clean) != ""

	// Candidate files, in priority order. vite-react-ssg emits flat files
	// (e.g. "cours-echecs-adultes-geneve.html"); the "/index.html" form covers
	// any nested-folder output too.
	var candidates []string
	switch {
	case clean == "/" || clean == ".":
		candidates = []string{"index.html"}
	case hasExt:
		candidates = []string{strings.TrimPrefix(clean, "/")}
	default:
		trimmed := strings.TrimPrefix(clean, "/")
		candidates = []string{trimmed + ".html", trimmed + "/index.html"}
	}

	for _, candidate := range candidates {
		if s.tryServe(w, r, candidate) {
			return
		}
	}
	// SPA-style fallback would hide broken links; instead serve the real 404
	// page so crawlers get a proper status.
	w.WriteHeader(http.StatusNotFound)
	if !s.tryServe(w, r, "404.html") {
		_, _ = io.WriteString(w, "404 not found")
	}
}

func (s *Server) tryServe(w http.ResponseWriter, r *http.Request, name string) bool {
	f, err := s.static.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil || info.IsDir() {
		return false
	}

	setContentType(w, name)
	setCacheHeaders(w, name)

	rs, ok := f.(io.ReadSeeker)
	if ok {
		http.ServeContent(w, r, name, modTime(info), rs)
		return true
	}
	// Embedded files may not be seekable depending on the io/fs impl.
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return true
	}
	_, _ = io.Copy(w, f)
	return true
}

func setContentType(w http.ResponseWriter, name string) {
	if ct := mime.TypeByExtension(path.Ext(name)); ct != "" {
		w.Header().Set("Content-Type", ct)
	}
}

func setCacheHeaders(w http.ResponseWriter, name string) {
	switch {
	case strings.HasPrefix(name, "assets/"):
		// Content-hashed by Vite → safe to cache forever.
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	case strings.HasSuffix(name, ".html"):
		// Always revalidate HTML so edits publish immediately.
		w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
	default:
		w.Header().Set("Cache-Control", "public, max-age=3600")
	}
}

func modTime(info fs.FileInfo) time.Time {
	if t := info.ModTime(); !t.IsZero() {
		return t
	}
	return time.Time{}
}
