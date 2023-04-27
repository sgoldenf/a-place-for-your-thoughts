package templates

import (
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/sgoldenf/a-place-for-your-thoughts/internal/models"
)

type TemplateData struct {
	Post            *models.Post
	TextMD          template.HTML
	Posts           []*models.Post
	PrevPage        int
	NextPage        int
	Form            any
	Popup           string
	IsAuthenticated bool
	CSRFToken       string
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./resources/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	fmt.Println(pages)
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.ParseFiles("./resources/html/pages/base.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./resources/html/components/*.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}
