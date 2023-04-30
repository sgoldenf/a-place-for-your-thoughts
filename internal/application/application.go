package application

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
<<<<<<< HEAD
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/adapters/post"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/adapters/user"
=======
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/models"
>>>>>>> 1a178fa (refactor: moved application logic to internal/application package && tested ping handler && tested testSecureHeaders middlware)
)

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
<<<<<<< HEAD
	Posts          post.PostModelInterface
	Users          user.UserModelInterface
=======
	Posts          *models.PostModel
	Users          *models.UserModel
>>>>>>> 1a178fa (refactor: moved application logic to internal/application package && tested ping handler && tested testSecureHeaders middlware)
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}
