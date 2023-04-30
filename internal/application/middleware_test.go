package application

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	secureHeaders(next).ServeHTTP(rr, r)
	rs := rr.Result()
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)
	expectedValue = "origin-when-cross-origin"
	equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)
	expectedValue = "nosniff"
	equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)
	expectedValue = "deny"
	equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)
	expectedValue = "0"
	equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)
	equal(t, rs.StatusCode, http.StatusOK)
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	equal(t, string(body), "OK")
}
