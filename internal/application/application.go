package application

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/models"
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	Posts          *models.PostModel
	Users          *models.UserModel
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}
