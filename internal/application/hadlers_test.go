package application

import (
	"net/http"
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
