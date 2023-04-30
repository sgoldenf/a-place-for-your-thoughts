package application

import (
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.Routes())
	defer ts.Close()
	code, _, body := ts.get(t, "/ping")
	equal(t, code, http.StatusOK)
	equal(t, string(body), "OK")
}

func TestViewPost(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.Routes())
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
			code, _, body := ts.get(t, test.urlPath)
			equal(t, code, test.wantCode)
			if test.wantTitle != "" {
				stringContains(t, string(body), test.wantTitle)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.Routes())
	defer ts.Close()
	_, _, body := ts.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)
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
			code, _, body := ts.postForm(t, "/user/signup", form)
			equal(t, code, test.wantCode)
			if test.wantFormTag != "" {
				stringContains(t, body, test.wantFormTag)
			}
		})
	}
}
