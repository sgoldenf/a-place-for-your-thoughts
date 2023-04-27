package application

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sgoldenf/a-place-for-your-thoughts/internal/testutils"
)

func TestPing(t *testing.T) {
	app := &Application{
		ErrorLog: log.New(io.Discard, "", 0),
		InfoLog:  log.New(io.Discard, "", 0),
	}
	ts := httptest.NewTLSServer(app.Routes())
	defer ts.Close()
	rs, err := ts.Client().Get(ts.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}
	testutils.Equal(t, rs.StatusCode, http.StatusOK)
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	testutils.Equal(t, string(body), "OK")
}
