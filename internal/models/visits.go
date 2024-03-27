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
	Duration    int
}

type CreateVisitDto struct {
	PetId       int       `json:"petId"`
	VetId       int       `json:"vetId"`
	Appointment time.Time `json:"appointment"`
	VisitReason string    `json:"visitReason"`
	Duration    int       `json:"duration"`
}

type VisitDetailDto struct {
	Id           int       `json:"id"`
	PetId        int       `json:"petId"`
	PetName      string    `json:"petName"`
	PetBirthdate time.Time `json:"petBirthdate"`
	PetType      string    `json:"petType"`
	VetId        int       `json:"vetId"`
	VetFirstName string    `json:"vetFirstName"`
	VetLastName  string    `json:"vetLastName"`
	Appointment  time.Time `json:"appointment"`
	VisitReason  string    `json:"visitReason"`
	Duration     int       `json:"duration"`
}

type VisitModelInterface interface {
	Create(petId, vetId int, appt time.Time, visitReason string, duration int) error
	GetByVetId(vetId int) ([]VisitDetailDto, error)
	GetById(visitId int) (VisitDetailDto, error)
	Remove(visitId int) error
}

type VisitModel struct {
	DB *sql.DB
}

func (model *VisitModel) Create(petId, vetId int, appt time.Time, visitReason string, duration int) error {

	stmt := `
		INSERT INTO visits (petId, vetId, appointment, created, visitReason, duration)
		VALUES (?, ?, ?, UTC_TIMESTAMP(), ?, ?)
	`

	_, err := model.DB.Exec(stmt, petId, vetId, appt, visitReason, duration)
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
			pet.birthdate,
			petType.name,
			vet.id,
			vet.firstName,
			vet.lastName,
			visit.appointment,
			visit.visitReason,
			visit.duration
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
		var visit VisitDetailDto
		err := rows.Scan(&visit.Id, &visit.PetId, &visit.PetName, &visit.PetBirthdate,
			&visit.PetType, &visit.VetId, &visit.VetFirstName, &visit.VetLastName,
			&visit.Appointment, &visit.VisitReason, &visit.Duration)

		if err != nil {
			return visits, err
		}

		visits = append(visits, visit)
	}

	return visits, nil
}

func (model *VisitModel) GetById(visitId int) (VisitDetailDto, error) {
	visit := VisitDetailDto{}

	stmt := `
		SELECT
			visit.id,
			pet.id,
			pet.name,
			pet.birthdate,
			petType.name,
			vet.id,
			vet.firstName,
			vet.lastName,
			visit.appointment,
			visit.visitReason,
			visit.duration
		FROM
			visits visit
			INNER JOIN vets vet on vet.id = visit.vetId
			INNER JOIN pets pet on pet.id = visit.petId
			INNER JOIN petTypes petType on petType.id = pet.petTypeId
		WHERE
			visit.id = ?;
	`

	err := model.DB.QueryRow(stmt, visitId).Scan(&visit.Id, &visit.PetId, &visit.PetName,
		&visit.PetBirthdate, &visit.PetType, &visit.VetId, &visit.VetFirstName,
		&visit.VetLastName, &visit.Appointment, &visit.VisitReason, &visit.Duration)

	if err != nil {
		return visit, err
	}

	return visit, nil
}

func (model *VisitModel) Remove(visitId int) error {

	stmt := `
		DELETE FROM visits
		WHERE id = ?
	`

	_, err := model.DB.Exec(stmt, visitId)
	if err != nil {
		return err
	}

	return nil
}
