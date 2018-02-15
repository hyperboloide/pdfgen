package main

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
)

type ServerEmulator struct {
	Data map[string]interface{}
	Tmpl *Template
	ts   *httptest.Server
}

func (s *ServerEmulator) Close() {
	s.ts.Close()
}

func (s *ServerEmulator) BaseURL() string {
	return s.ts.URL
}

func (s *ServerEmulator) MainHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.Tmpl.Index.ExecuteWriterUnbuffered(s.Data, w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *ServerEmulator) FooterHandler(w http.ResponseWriter, r *http.Request) {
	if s.Tmpl.Footer == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else if err := s.Tmpl.Footer.ExecuteWriterUnbuffered(s.Data, w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

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
