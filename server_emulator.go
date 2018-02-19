package main

import (
	"net/http"
	"net/http/httptest"

	"github.com/flosch/pongo2"
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

// HandleTemplate will generate the HTML from the template and handle
// the http request.
func (s *ServerEmulator) HandleTemplate(tmpl *pongo2.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if tmpl == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else if err := tmpl.ExecuteWriterUnbuffered(s.Data, w); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

// NewServerEmulator creates a new temporary server for the lifetime of the request.
func NewServerEmulator(d map[string]interface{}, t *Template) *ServerEmulator {
	s := &ServerEmulator{
		Data: d,
		Tmpl: t,
	}
	r := chi.NewRouter()
	r.HandleFunc("/main", s.HandleTemplate(s.Tmpl.Index))
	r.HandleFunc("/header", s.HandleTemplate(s.Tmpl.Header))
	r.HandleFunc("/footer", s.HandleTemplate(s.Tmpl.Footer))
	r.Mount("/", http.FileServer(http.Dir(s.Tmpl.RootDir)))
	s.ts = httptest.NewServer(r)
	return s
}
