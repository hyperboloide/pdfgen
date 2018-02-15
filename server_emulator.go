package main

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
)

// ServerEmulator is a temporary http server only for wkhtmltopdf to download
// the html build from the templates
type ServerEmulator struct {
	Data map[string]interface{}
	Tmpl *Template
	ts   *httptest.Server
}

// Close the temporary server
func (s *ServerEmulator) Close() {
	s.ts.Close()
}

// BaseURL returns the temporary server url
func (s *ServerEmulator) BaseURL() string {
	return s.ts.URL
}

// MainHandler will build the main template (index.html) and render it
func (s *ServerEmulator) MainHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.Tmpl.Index.ExecuteWriterUnbuffered(s.Data, w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// FooterHandler will build the footer template (footer.html) and render it
func (s *ServerEmulator) FooterHandler(w http.ResponseWriter, r *http.Request) {
	if s.Tmpl.Footer == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else if err := s.Tmpl.Footer.ExecuteWriterUnbuffered(s.Data, w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// NewServerEmulator creates a new temporary server for the lifetime of the request.
func NewServerEmulator(d map[string]interface{}, t *Template) *ServerEmulator {
	s := &ServerEmulator{
		Data: d,
		Tmpl: t,
	}
	r := chi.NewRouter()
	r.HandleFunc("/main", s.MainHandler)
	r.HandleFunc("/footer", s.FooterHandler)
	r.Mount("/", http.FileServer(http.Dir(s.Tmpl.RootDir)))
	s.ts = httptest.NewServer(r)
	return s
}
