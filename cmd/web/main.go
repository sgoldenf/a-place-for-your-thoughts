package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/application"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/models"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/templates"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dbURL := flag.String("dbURL", "postgres://sgoldenf:sgoldenf@localhost:5432/blog", "PostgresSQL database URL")
	cert := flag.String("tls-cert", "./tls/cert.pem", "TLS public key")
	key := flag.String("tls-key", "./tls/key.pem", "TLS private key")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	pool, err := dbConn(*dbURL)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer pool.Close()
	cache, err := templates.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	app := &application.Application{
		ErrorLog:       errorLog,
		InfoLog:        infoLog,
		Posts:          &models.PostModel{Pool: pool},
		Users:          &models.UserModel{Pool: pool},
		TemplateCache:  cache,
		FormDecoder:    formDecoder,
		SessionManager: sessionManager,
	}
	tlsConfig := &tls.Config{CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}}
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS(*cert, *key)
	errorLog.Fatal(err)
}

func dbConn(dbURL string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return conn, err
}
