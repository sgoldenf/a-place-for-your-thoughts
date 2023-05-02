package application

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/models/mocks"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/templates"
	testutils "github.com/sgoldenf/a-place-for-your-thoughts/internal/test_utils"
)

const templateTestPath = "../../resources/html"

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

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()
	code, _, body := ts.Get(t, "/ping")
	testutils.Equal(t, code, http.StatusOK)
	testutils.Equal(t, string(body), "OK")
}

func TestViewPost(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()
	tests := []struct {
		name      string
		urlPath   string
		wantCode  int
		wantTitle string
	}{
		{
			name:      "Valid ID",
			urlPath:   "/post/view/1",
			wantCode:  http.StatusOK,
			wantTitle: "Title 1",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/post/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/post/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Float ID",
			urlPath:  "/post/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/post/view/text",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/post/view/",
			wantCode: http.StatusNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, _, body := ts.Get(t, test.urlPath)
			testutils.Equal(t, code, test.wantCode)
			if test.wantTitle != "" {
				testutils.StringContains(t, string(body), test.wantTitle)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()
	_, _, body := ts.Get(t, "/user/signup")
	validCSRFToken := testutils.ExtractCSRFToken(t, body)
	const (
		validName     = "Sgoldenf"
		validPassword = "validPassword"
		validEmail    = "sgoldenf@example.com"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)
	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{
			name:         "Valid Submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    "wrongToken",
			wantCode:     http.StatusBadRequest,
		},
		{
			name:         "Empty Name",
			userName:     "",
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty Email",
			userName:     validName,
			userEmail:    "",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Empty Password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Invalid Email",
			userName:     validName,
			userEmail:    "invalidEmail",
			userPassword: validEmail,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Short Password",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: "pass",
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:         "Duplicate Email",
			userName:     validName,
			userEmail:    "duplicate@example.com",
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", test.userName)
			form.Add("email", test.userEmail)
			form.Add("password", test.userPassword)
			form.Add("csrf_token", test.csrfToken)
			code, _, body := ts.PostForm(t, "/user/signup", form)
			testutils.Equal(t, code, test.wantCode)
			if test.wantFormTag != "" {
				testutils.StringContains(t, body, test.wantFormTag)
			}
		})
	}
}

func TestCreatePostForm(t *testing.T) {
	app := newTestApplication(t)
	ts := testutils.NewTestServer(t, app.Routes())
	defer ts.Close()
	t.Run("Unauthenticated", func(f *testing.T) {
		code, header, _ := ts.Get(t, "/post/create")
		testutils.Equal(t, code, http.StatusSeeOther)
		testutils.Equal(t, header.Get("Location"), "/user/login")
	})
	t.Run("Unauthenticated", func(t *testing.T) {
		_, _, body := ts.Get(t, "/user/signup")
		csrfToken := testutils.ExtractCSRFToken(t, body)
		form := url.Values{}
		form.Add("email", "login@example.com")
		form.Add("password", "password")
		form.Add("csrf_token", csrfToken)
		ts.PostForm(t, "/user/login", form)
		code, _, body := ts.Get(t, "/post/create")
		testutils.Equal(t, code, http.StatusOK)
		testutils.StringContains(t, body, "<form action='/post/create' method='POST'>")
	})
}
