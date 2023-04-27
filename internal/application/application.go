package application

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/adapters/post"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/adapters/user"
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	Posts          post.PostModelInterface
	Users          user.UserModelInterface
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}
