package application

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/models/mocks"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/templates"
)

type testServer struct {
	*httptest.Server
}

func newTestApplication(t *testing.T) *Application {
	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		t.Fatal(err)
	}
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	return &Application{
		ErrorLog:       log.New(io.Discard, "", 0),
		InfoLog:        log.New(io.Discard, "", 0),
		Posts:          &mocks.PostModel{},
		Users:          &mocks.UserModel{},
		TemplateCache:  templateCache,
		FormDecoder:    form.NewDecoder(),
		SessionManager: sessionManager,
	}
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	ts.Client().Jar = jar
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}

func equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("got: %v; expected: %v", actual, expected)
	}
}

func stringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()
	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}
