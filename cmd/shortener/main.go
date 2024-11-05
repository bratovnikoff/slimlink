package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/sqids/sqids-go"
)

var (
	urls = make(map[string]string)
	s, _ = sqids.New(sqids.Options{
		Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.+",
		MinLength: 8,
	})
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	log.Fatal(http.ListenAndServe(`:8080`, mux))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		encodeURLHandler(w, r)
	} else if r.Method == http.MethodGet {
		resolveURLHandler(w, r)
	} else {
		http.Error(w, "This request is not allowed.", http.StatusBadRequest)
	}
}

func encodeURLHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	if r.URL.Path != "/" || r.Method != http.MethodPost || len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, _ := s.Encode([]uint64{uint64(len(urls))})

	urls[id] = string(body)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("http://" + r.Host + "/" + id))
}

func resolveURLHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	url, ok := urls[id]
	if r.Method != http.MethodGet || !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
