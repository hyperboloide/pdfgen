package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
)

type Fake struct {
	Addr string
}

func GenHandler(tmpl string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		s, ok := Sessions.Get(vars["session"])
		if !ok {
			statError(w, http.StatusNotFound)
			return
		}
		t := Templates[vars["template"]]
		if err := t.template.ExecuteTemplate(w, tmpl, s); err != nil {
			statError(w, http.StatusInternalServerError)
			return
		}
	}
}

func StartFake() string {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		log.Fatal("could not open loopback connection")
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal("could not use loopback connection")
	}
	port := l.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d", port)

	r := mux.NewRouter()
	r.HandleFunc("/{template}/{session}/index", GenHandler("index"))
	r.HandleFunc("/{template}/{session}/footer", GenHandler("footer"))
	r.HandleFunc("/{template}/{session}/header", GenHandler("header"))
	for k, v := range Templates {
		path := fmt.Sprintf("/{%s}/{session}/", k)
		r.PathPrefix(path).Handler(http.FileServer(http.Dir(v.rootDir)))
	}
	go func() {
		if err := http.Serve(l, r); err != nil {
			log.Fatal(err)
		}
	}()
	return url
}
