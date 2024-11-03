package main

import (
	"io"
	"net/http"

	"github.com/sqids/sqids-go"
)

var (
	urls map[string]string
	s    *sqids.Sqids
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	urls = make(map[string]string)
	s, _ = sqids.New(sqids.Options{
		Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.+",
		MinLength: 8,
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/", checkMethod)

	return http.ListenAndServe(`:8080`, mux)
}

func checkMethod(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		encodeURL(w, r)
	} else if r.Method == http.MethodGet {
		resolveURL(w, r)
	}
}

func encodeURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if r.URL.Path != "/" || r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := s.Encode([]uint64{uint64(len(urls))})
	if err != nil {
		panic(err)
	}

	urls[id] = string(body)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("http://" + r.Host + "/" + id))
}

func resolveURL(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/"):]
	url, ok := urls[id]
	if id == "" || !ok || r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
