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

type VisitDetailDto struct {
	Id           int
	PetId        int
	PetName      string
	PetType      string
	VetId        int
	VetFirstName string
	VetLastName  string
	Appointment  time.Time
	VisitReason  string
}

type VisitModelInterface interface {
	Create(petId, vetId int, appt time.Time, visitReason string) error
	GetByVetId(vetId int) ([]VisitDetailDto, error)
}

type VisitModel struct {
	DB *sql.DB
}

func (model *VisitModel) Create(petId, vetId int, appt time.Time, visitReason string) error {

	stmt := `
		INSERT INTO visits (petId, vetId, appointment, created, visitReason)
		VALUES (?, ?, ?, UTC_TIMESTAMP(), ?)
	`

	_, err := model.DB.Exec(stmt, petId, vetId, appt, visitReason)
	if err != nil {
		return err
	}

	return nil
}

func (model *VisitModel) GetByVetId(vetId int) ([]VisitDetailDto, error) {
	visits := []VisitDetailDto{}

	stmt := `
		SELECT
			visit.id,
			pet.id,
			pet.name,
			petType.name,
			vet.id,
			vet.firstName,
			vet.lastName,
			visit.appointment,
			visit.visitReason
		FROM
			visits visit
			INNER JOIN vets vet on vet.id = visit.vetId
			INNER JOIN pets pet on pet.id = visit.petId
			INNER JOIN petTypes petType on petType.id = pet.petTypeId
		WHERE
			vet.id = ?;
	`

	rows, err := model.DB.Query(stmt, vetId)
	if err != nil {
		return visits, err
	}

	for rows.Next() {
		var visitDetail VisitDetailDto
		err := rows.Scan(&visitDetail.Id, &visitDetail.PetId, &visitDetail.PetName,
			&visitDetail.PetType, &visitDetail.VetId, &visitDetail.VetFirstName,
			&visitDetail.VetLastName, &visitDetail.Appointment, &visitDetail.VisitReason)

		if err != nil {
			return visits, err
		}

		visits = append(visits, visitDetail)
	}

	return visits, nil
}
