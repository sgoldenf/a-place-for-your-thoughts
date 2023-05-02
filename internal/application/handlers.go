package application

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/models"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/validator"
	"github.com/yuin/goldmark"
)

type postCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	validator.Validator `form:"-"`
}

type userSignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	count, err := app.Posts.GetPostsCount()
	if err != nil {
		app.serverError(w, err)
		return
	}
	page, err := getPage(r, count)
	if err != nil {
		app.notFound(w)
		return
	}
	posts, err := app.Posts.GetPostsList(page)
	if err != nil {
		app.serverError(w, err)
		return
	}
	nextPage := page + 1
	if page*10+1 > count {
		nextPage = 0
	}
	data := app.NewTemplateData(r)
	data.Posts = posts
	data.PrevPage = page - 1
	data.NextPage = nextPage
	app.render(
		w, http.StatusOK, "home.tmpl", data)
}

func getPage(r *http.Request, postCount int) (int, error) {
	page := chi.URLParam(r, "page")
	if page == "" {
		page = "1"
	}
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		return 0, err
	}
	if (pageNumber > 1 && (pageNumber-1)*10 >= postCount) || pageNumber < 1 {
		return 0, errors.New("not found")
	}
	return pageNumber, nil
}

func (app *Application) createPostForm(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = postCreateForm{}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *Application) CreatePost(w http.ResponseWriter, r *http.Request) {
	var form postCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 255), "title", "Max length for this field is 255 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	if !form.Valid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
	id, err := app.Posts.CreatePost(form.Title, form.Content)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.SessionManager.Put(r.Context(), "popup", "Post successfuly created!")
	http.Redirect(w, r, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

func (app *Application) viewPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	post, err := app.Posts.ReadPost(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.NewTemplateData(r)
	data.Post = post
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(post.Text), &buf); err == nil {
		data.TextMD = template.HTML(buf.Bytes())
	}
	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = userSignUpForm{}
	app.render(w, http.StatusOK, "signup.tmpl", data)
}

func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignUpForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Invalid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	if !form.Valid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}
	err = app.Users.CreateUser(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is alredy in use")
			data := app.NewTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.SessionManager.Put(r.Context(), "popup", "Your signup was successful. Please log in")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Invalid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	if !form.Valid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}
	id, err := app.Users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.NewTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	err = app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.SessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/post/create", http.StatusSeeOther)
}

func (app *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
	}
	app.SessionManager.Remove(r.Context(), "authenticatedUserID")
	app.SessionManager.Put(r.Context(), "popup", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
