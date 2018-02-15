package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
	"github.com/streamrail/concurrent-map"
)

var (
	Sessions = cmap.New()
)

func statError(w http.ResponseWriter, status int) {
	msg := fmt.Sprintf("%d - %s", status, http.StatusText(status))
	http.Error(w, msg, status)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "template")
	var data map[string]interface{}

	if r.Method != "POST" && r.Method != "PUT" {
		statError(w, http.StatusMethodNotAllowed)
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

func main() {
	configRead()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.DefaultCompress)
	r.HandleFunc("/{template}", Handler)

	addr := fmt.Sprintf("%s:%d",
		viper.GetString("addr"),
		viper.GetInt("port"),
	)
	log.Printf("accepting connections on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
