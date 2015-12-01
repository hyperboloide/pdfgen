package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/streamrail/concurrent-map"
	"log"
	"math/rand"
	"net/http"
	"os"
)

var (
	Sessions = cmap.New()
	FakeUrl  string
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
	var data interface{}
	if err := decoder.Decode(data); err != nil {
		statError(w, http.StatusBadRequest)
		return
	}
	id := fmt.Sprintf("%d", rand.Int())
	Sessions.Set(id, data)
	defer Sessions.Remove(id)
	w.Header().Set("Content-Type", "application/pdf")
	if err := tmpl.Gen(id, w); err != nil {
		statError(w, http.StatusInternalServerError)
		return
	}
}

func main() {
	configRead()

	nb := len(Templates)
	switch nb {
	case 0:
		fmt.Println("No templates found, exiting.")
		return
	case 1:
		fmt.Println("1 template found:")
	default:
		fmt.Printf("%d templates found:\n", nb)
	}
	for k, _ := range Templates {
		fmt.Printf("  - %s\n", k)
	}

	FakeUrl = StartFake()

	r := mux.NewRouter()
	r.HandleFunc("/{template}", Handler)
	mddw := handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(r))
	fmt.Printf("accepting connections on %s:%s\n", Addr, Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", Addr, Port), mddw); err != nil {
		log.Fatal(err)
	}
}
