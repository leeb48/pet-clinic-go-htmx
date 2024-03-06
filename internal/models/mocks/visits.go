package mocks

import (
	"time"

	"pet-clinic.bonglee.com/internal/models"
)

type VisitModel struct {
}

func (model *VisitModel) Create(petId, vetId int, appt time.Time, visitReason string, duration int) error {

	return nil
}

func (model *VisitModel) GetByVetId(vetId int) ([]models.VisitDetailDto, error) {
	visits := []models.VisitDetailDto{}

	return visits, nil
}
