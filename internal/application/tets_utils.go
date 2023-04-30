package application

import (
	"bytes"
	"context"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/models/mocks"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/templates"
)

type testServer struct {
	*httptest.Server
}

const templateTestPath = "../../resources/html"

var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

func newTestApplication(t *testing.T) *Application {
	templateCache, err := templates.NewTemplateCache(templateTestPath)
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

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
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

func extractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
	return html.UnescapeString(string(matches[1]))
}

func newTestDB(t *testing.T) *pgxpool.Pool {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatal(err)
	}
	testDB := os.Getenv("TEST_DB")
	user := os.Getenv("TEST_DB_USER")
	password := os.Getenv("TEST_DB_PASSWORD")
	port := os.Getenv("TEST_DB_PORT")
	dbURL := "postgres://" + user + ":" + password + "@localhost:" + port + "/" + testDB
	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatal(err)
	}
	script, err := os.ReadFile("../models/testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(context.Background(), string(script))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		script, err := os.ReadFile("../models/testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(context.Background(), string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})
	return db
}
