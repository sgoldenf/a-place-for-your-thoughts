package application

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
<<<<<<< HEAD
=======

	"github.com/sgoldenf/a-place-for-your-thoughts/internal/testutils"
>>>>>>> 1a178fa (refactor: moved application logic to internal/application package && tested ping handler && tested testSecureHeaders middlware)
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
<<<<<<< HEAD
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
=======
	testutils.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)
	expectedValue = "origin-when-cross-origin"
	testutils.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)
	expectedValue = "nosniff"
	testutils.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)
	expectedValue = "deny"
	testutils.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)
	expectedValue = "0"
	testutils.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)
	testutils.Equal(t, rs.StatusCode, http.StatusOK)
>>>>>>> 1a178fa (refactor: moved application logic to internal/application package && tested ping handler && tested testSecureHeaders middlware)
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
<<<<<<< HEAD
	equal(t, string(body), "OK")
=======
	testutils.Equal(t, string(body), "OK")
>>>>>>> 1a178fa (refactor: moved application logic to internal/application package && tested ping handler && tested testSecureHeaders middlware)
}
