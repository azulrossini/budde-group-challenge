package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:dist
var distFS embed.FS

type Handler struct {
	fs http.FileSystem
}

func New() (*Handler, error) {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil, err
	}
	return &Handler{fs: http.FS(sub)}, nil
}

// ServeHTTP serves a requested file if it exists in the built SPA, and
// falls back to index.html otherwise so client-side navigation and
// hard-refreshes both work.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := h.fs.Open(r.URL.Path)
	if err != nil {
		fallback := new(http.Request)
		*fallback = *r
		fallback.URL.Path = "/"
		http.FileServer(h.fs).ServeHTTP(w, fallback)
		return
	}
	f.Close()

	http.FileServer(h.fs).ServeHTTP(w, r)
}
