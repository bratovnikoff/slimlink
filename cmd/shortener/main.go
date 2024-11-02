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
		decodeURL(w, r)
	}
}

func encodeURL(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	if req.URL.Path != "/" || contentType != "text/plain" || len(body) == 0 {
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

	_, _ = w.Write([]byte("http://localhost:8080/" + id + "\r\n"))
}

func decodeURL(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[len("/"):]
	_, ok := urls[id]
	if id == "" || !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addr := urls[id]

	w.Header().Set("Location", addr)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
