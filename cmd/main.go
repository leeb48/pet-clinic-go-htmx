package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/caarlos0/env"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"pet-clinic.bonglee.com/internal/app"
	"pet-clinic.bonglee.com/internal/constants"
	"pet-clinic.bonglee.com/internal/handlers"
	"pet-clinic.bonglee.com/internal/models"
)

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

	cfg := app.Config{}
	if err := env.Parse(&cfg); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	dsn := cfg.LOCAL_DSN
	if cfg.ENV == constants.PROD {
		dsn = cfg.PROD_DSN
	}
	db, err := openDB(dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := app.CreateTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	session := scs.New()
	session.Lifetime = time.Duration(cfg.SessionDuration) * time.Hour

	app := app.App{
		Logger:        logger,
		TemplateCache: templateCache,
		Cfg:           cfg,
		Session:       session,
		Owners:        &models.OwnerModel{DB: db},
		PetTypes:      &models.PetTypeModel{DB: db},
		Pets:          &models.PetModel{DB: db},
		Vets:          &models.VetModel{DB: db},
		Visits:        &models.VisitModel{DB: db},
	}

	server := http.Server{
		Addr:         cfg.Addr,
		Handler:      handlers.Routes(&app),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.Logger.Info("starting server", slog.String("addr", app.Cfg.Addr))

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

	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		db.Close()
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)

	if err != nil {
		db.Close()
		return nil, err
	}

	m.Steps(11)

	return db, nil
}
