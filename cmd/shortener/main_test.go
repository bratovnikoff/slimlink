package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_encodeURLHandler(t *testing.T) {

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
			body := strings.NewReader("https://www.youtube.com/")
			r := httptest.NewRequest(tt.method, "/", body)
			w := httptest.NewRecorder()

			encodeURLHandler(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func Test_resolveURLHadler(t *testing.T) {
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
			urls["bbbbbbbb"] = "ffff"

			r := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()

			resolveURLHandler(w, r)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}
