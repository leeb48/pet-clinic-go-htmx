package main

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger        *slog.Logger
	templateCache map[string]*template.Template
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	templateCache, err := createTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := application{
		logger:        logger,
		templateCache: templateCache,
	}

	err = http.ListenAndServe(":7771", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
