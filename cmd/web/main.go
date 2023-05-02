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
	"github.com/joho/godotenv"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/application"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/models"
	"github.com/sgoldenf/a-place-for-your-thoughts/internal/templates"
)

const templateBasePath = "./resources/html"

var (
	addr  *string
	dbURL *string
	cert  *string
	key   *string
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("WARNING: No .env file found")
	}
	addr = flag.String("addr", ":4000", "HTTP network address")
	dbName := os.Getenv("APP_DB")
	user := os.Getenv("APP_DB_USER")
	password := os.Getenv("APP_DB_PASSWORD")
	dbURL = flag.String(
		"dbURL",
		"postgres://"+user+":"+password+"@localhost:5432/"+dbName,
		"PostgresSQL database URL",
	)
	cert = flag.String("tls-cert", os.Getenv("TLS_CERT"), "TLS public key")
	key = flag.String("tls-key", os.Getenv("TLS_KEY"), "TLS private key")
	flag.Parse()
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	pool, err := dbConn()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer pool.Close()
	cache, err := templates.NewTemplateCache(templateBasePath)
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

func dbConn() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), *dbURL)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return conn, err
}
