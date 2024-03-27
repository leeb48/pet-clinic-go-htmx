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

func (model *VisitModel) GetById(visitId int) (models.VisitDetailDto, error) {
	visit := models.VisitDetailDto{}

	return visit, nil
}

func (model *VisitModel) Remove(visitId int) error {

	return nil
}
