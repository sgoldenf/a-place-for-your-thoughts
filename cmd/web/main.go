package main

import (
	"context"
	"crypto/tls"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/models"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	posts          *models.PostModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

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
	cache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		posts:          &models.PostModel{Pool: pool},
		users:          &models.UserModel{Pool: pool},
		templateCache:  cache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
	tlsConfig := &tls.Config{CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}}
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
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
