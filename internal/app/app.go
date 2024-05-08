package app

import (
	"html/template"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"pet-clinic.bonglee.com/internal/models"
)

type App struct {
	Logger        *slog.Logger
	TemplateCache map[string]*template.Template
	Cfg           Config
	Session       *scs.SessionManager
	Owners        models.OwnerModelInterface
	PetTypes      models.PetTypeModelInterface
	Pets          models.PetModelInterface
	Vets          models.VetModelInterface
	Visits        models.VisitModelInterface
}

type Config struct {
	ENV             string `env:"ENV"`
	Addr            string `env:"ADDR"`
	LOCAL_DSN       string `env:"LOCAL_DSN"`
	PROD_DSN        string `env:"PROD_DSN"`
	SessionDuration int    `env:"SESSION_DURATION"`
	DefaultPetType  int    `env:"DEFAULT_PET_TYPE"`
}
