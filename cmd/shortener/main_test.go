package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sqids/sqids-go"
	"github.com/stretchr/testify/assert"
)

func Test_encodeURL(t *testing.T) {

	tests := []struct {
		method       string
		expectedCode int
	}{
		{method: http.MethodGet, expectedCode: http.StatusBadRequest},
		{method: http.MethodPut, expectedCode: http.StatusBadRequest},
		{method: http.MethodDelete, expectedCode: http.StatusBadRequest},
		{method: http.MethodPost, expectedCode: http.StatusCreated},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			urls = make(map[string]string)
			s, _ = sqids.New(sqids.Options{
				Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.+",
				MinLength: 8,
			})

			r := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			encodeURL(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func Test_resolveURL(t *testing.T) {
	tests := []struct {
		method       string
		request      string
		expectedCode int
	}{
		{method: http.MethodPost, request: "/eeeeeeee", expectedCode: http.StatusBadRequest},
		{method: http.MethodPut, request: "/eeeeeeee", expectedCode: http.StatusBadRequest},
		{method: http.MethodDelete, request: "/eeeeeeee", expectedCode: http.StatusBadRequest},
		{method: http.MethodGet, request: "/bbbbbbbb", expectedCode: http.StatusTemporaryRedirect},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			urls = map[string]string{"bbbbbbbb": "ffff"}
			r := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()

			resolveURL(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
