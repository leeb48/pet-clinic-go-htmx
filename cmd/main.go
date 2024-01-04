package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"pet-clinic.bonglee.com/internal/models"
)

type application struct {
	logger        *slog.Logger
	templateCache map[string]*template.Template
	owners        models.OwnerModelInterface
	petTypes      models.PetTypeModelInterface
}

type config struct {
	addr string
	dsn  string
}

func main() {

	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":7771", "HTTP network address")
	flag.StringVar(&cfg.dsn, "dsn", "app:1234@/petClinic?parseTime=true&multiStatements=true", "MySQL data source name")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := createTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := application{
		logger:        logger,
		templateCache: templateCache,
		owners:        &models.OwnerModel{DB: db},
		petTypes:      &models.PetTypeModel{DB: db},
	}

	server := http.Server{
		Addr:         cfg.addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
