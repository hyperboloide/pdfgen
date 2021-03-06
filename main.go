package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
)

func statError(w http.ResponseWriter, status int) {
	msg := fmt.Sprintf("%d - %s", status, http.StatusText(status))
	http.Error(w, msg, status)
}

// APIHandler is the public facing http handler, it will respond only to POST or
// PUT requests that match the template name (ex: /my_template)
func APIHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "template")
	var data map[string]interface{}

	if r.Method != "POST" && r.Method != "PUT" {
		statError(w, http.StatusMethodNotAllowed)

	} else if r.Header.Get("Content-type") != "application/json" {
		statError(w, http.StatusBadRequest)

	} else if tmpl, exists := Templates[name]; !exists {
		statError(w, http.StatusNotFound)

	} else if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		statError(w, http.StatusBadRequest)

	} else {
		w.Header().Set("Content-Type", "application/pdf")
		srv := NewServerEmulator(data, tmpl)
		defer srv.Close()

		if err := tmpl.WritePDF(srv.BaseURL(), w); err != nil {
			log.Print(err)
			statError(w, http.StatusInternalServerError)
		}
	}

}

// Router builds the http router.
func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.DefaultCompress)
	r.HandleFunc("/{template}", APIHandler)
	return r
}

func main() {
	if err := ConfigRead(); err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("%s:%d", viper.GetString("addr"), viper.GetInt("port"))
	log.Printf("accepting connections on %s", addr)
	if err := http.ListenAndServe(addr, Router()); err != nil {
		log.Fatal(err)
	}
}
