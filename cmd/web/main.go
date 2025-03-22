// Filename:main.go
package main

import (
	"context"      // Provides context for DB operations to set timeouts
	"database/sql" // Provides functions for interacting with a database
	"flag"
	"html/template"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"                          // PostgreSQL driver (imported anonymously)
	"github.com/mickali02/templates/internal/data" // Importing custom package for feedback model
)

type application struct {
	logger        *slog.Logger
	addr          *string
	feedback      *data.FeedbackModel           // A model to interact with the feedback data in the DB
	templateCache map[string]*template.Template // Cache of HTML templates for faster rendering
}

func main() {
	addr := flag.String("addr", "", "HTTP network address")
	dsn := flag.String("dsn", "", "PostgreSQL DSN") // DSN is a connection string to the database

	flag.Parse()

	// Set up a structured logger for logging messages
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Open the connection pool to the PostgreSQL database
	db, err := openDB(*dsn)
	if err != nil { // Handle error if DB connection fails
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("database connection pool established") // Log success message

	// Load HTML templates into cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close() // Ensure the DB connection is closed when the program ends

	// Create an instance of the application, passing the logger, DB, etc.
	app := &application{
		logger:        logger,
		addr:          addr,
		feedback:      &data.FeedbackModel{DB: db},
		templateCache: templateCache,
	}

	// Start serving HTTP requests
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// openDB establishes a connection to the PostgreSQL database
func openDB(dsn string) (*sql.DB, error) {
	// open a connection pool to PostgreSQL using the DSN (Data Source Name)
	db, err := sql.Open("postgres", dsn)
	if err != nil { // If error occurs, return the error
		return nil, err
	}

	// Set a timeout context for DB operations to prevent long hangs (5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Ensure the context is cancelled to avoid memory leaks

	err = db.PingContext(ctx) // Ping the DB to check if it's reachable within the timeout
	if err != nil {           // If the DB is not reachable, close the connection and return error
		db.Close()
		return nil, err
	}

	// Return the open database connection pool
	return db, nil
}
