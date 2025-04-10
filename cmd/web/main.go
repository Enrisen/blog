package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"os"
	"time"

	"github.com/Enrisen/blog/internal/data"
	_ "github.com/lib/pq"
)

type application struct {
	logger        *slog.Logger
	addr          *string
	blog          *data.BlogModel
	users         *data.UserModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", "", "HTTP network address")
	blogDSN := flag.String("dsn", os.Getenv("it_blog_DB_DSN"), "Blog PostgreSQL DSN")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Open blog database connection
	blogDB, err := openDB(*blogDSN)
	if err != nil {
		logger.Error("failed to open blog database", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("blog database connection pool established")

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer blogDB.Close()

	app := &application{
		logger:        logger,
		addr:          addr,
		blog:          &data.BlogModel{DB: blogDB},
		users:         &data.UserModel{DB: blogDB},
		templateCache: templateCache,
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	// open a connection pool
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// set a context to ensure DB operations don't take too long
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// return the connection pool (sql.DB)
	return db, nil

}
