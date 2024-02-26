package models

import (
	"database/sql"
	"time"
)

type Visit struct {
	Id          int
	PetId       int
	VetId       int
	Appointment time.Time
	Created     time.Time
	VisitReason string
}

type CreateVisitDto struct {
	PetId       int       `json:"petId"`
	VetId       int       `json:"vetId"`
	Appointment time.Time `json:"appointment"`
	VisitReason string    `json:"visitReason"`
}

type VisitModelInterface interface {
}

type VisitModel struct {
	DB *sql.DB
}
