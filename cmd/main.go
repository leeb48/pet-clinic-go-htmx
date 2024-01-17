package main

import (
	"database/sql"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/caarlos0/env"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"pet-clinic.bonglee.com/internal/models"
)

type application struct {
	logger        *slog.Logger
	templateCache map[string]*template.Template
	cfg           config
	session       *scs.SessionManager
	owners        models.OwnerModelInterface
	petTypes      models.PetTypeModelInterface
	pets          models.PetModelInterface
}

type config struct {
	Addr            string `env:"ADDR"`
	DSN             string `env:"DSN"`
	SessionDuration int    `env:"SESSION_DURATION"`
	DefaultPetType  int    `env:"DEFAULT_PET_TYPE"`
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	db, err := openDB(cfg.DSN)
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

	session := scs.New()
	session.Lifetime = time.Duration(cfg.SessionDuration) * time.Hour

	app := application{
		logger:        logger,
		templateCache: templateCache,
		cfg:           cfg,
		session:       session,
		owners:        &models.OwnerModel{DB: db},
		petTypes:      &models.PetTypeModel{DB: db},
		pets:          &models.PetModel{DB: db},
	}

	server := http.Server{
		Addr:         cfg.Addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("starting server", slog.String("addr", app.cfg.Addr))

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
