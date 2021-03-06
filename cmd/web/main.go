package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.bimasenaputra/internal/models"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets models.SnippetModelInterface
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "web:password@/snippetbox?parseTime=true", "MySQL Data Source Name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application {
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	infoLog.Println("Starting server on", *addr)

	server := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}

	shutdownErr := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		errorLog.Println(s)
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		shutdownErr <- server.Shutdown(ctx)
	}()

	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}