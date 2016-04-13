package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	if r.Method != "POST" && r.Method != "PUT" {
		statError(w, http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	name := vars["template"]
	tmpl, exists := Templates[name]
	if exists == false {
		statError(w, http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		log.Print(err)
		statError(w, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	srv := NewServerEmulator(data, tmpl)
	defer srv.Close()

	if err := tmpl.WritePDF(srv.BaseURL(), w); err != nil {
		log.Print(err)
		statError(w, http.StatusInternalServerError)
	}
}

func main() {
	configRead()

	nb := len(Templates)
	switch nb {
	case 0:
		fmt.Println("No template found, exiting.")
		return
	case 1:
		fmt.Println("1 template found:")
	default:
		fmt.Printf("%d templates found:\n", nb)
	}
	for k, _ := range Templates {
		fmt.Printf("  - %s\n", k)
	}

	fmt.Printf("accepting connections on %s:%s\n", Addr, Port)

	r := mux.NewRouter()
	r.HandleFunc("/{template}", Handler)

	err := http.ListenAndServe(
		fmt.Sprintf("%s:%s", Addr, Port),
		handlers.LoggingHandler(
			os.Stdout,
			handlers.CompressHandler(r)))
	if err != nil {
		log.Fatal(err)
	}
}
