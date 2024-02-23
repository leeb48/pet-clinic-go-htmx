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
	PetId       int
	VetId       int
	Appointment time.Time
	VisitReason string
}

type VisitModelInterface interface {
}

type VisitModel struct {
	DB *sql.DB
}
