package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	}))
	fileServer := http.FileServer(http.Dir("./resources/static/"))
	router.Method(http.MethodGet, "/static/*", http.StripPrefix("/static", fileServer))
	router.Get("/", app.loadAndSaveSession(app.noSurf(app.authenticate(app.home))))
	router.Get("/page/{page}", app.loadAndSaveSession(app.noSurf(app.authenticate(app.home))))
	router.Get("/post/view/{id}", app.loadAndSaveSession(app.noSurf(app.authenticate(app.viewPost))))
	router.Get("/post/create", app.loadAndSaveSession(app.noSurf(app.authenticate(app.requireAuthentication(app.createPostForm)))))
	router.Post("/post/create", app.loadAndSaveSession(app.noSurf(app.authenticate(app.requireAuthentication(app.CreatePost)))))
	router.Get("/user/signup", app.loadAndSaveSession(app.noSurf(app.authenticate(app.userSignup))))
	router.Post("/user/signup", app.loadAndSaveSession(app.noSurf(app.authenticate(app.userSignupPost))))
	router.Get("/user/login", app.loadAndSaveSession(app.noSurf(app.authenticate(app.userLogin))))
	router.Post("/user/login", app.loadAndSaveSession(app.noSurf(app.authenticate(app.userLoginPost))))
	router.Post("/user/logout", app.loadAndSaveSession(app.noSurf(app.authenticate(app.requireAuthentication(app.userLogoutPost)))))
	return app.recoverPanic(app.logRequest(secureHeaders(router)))
}
