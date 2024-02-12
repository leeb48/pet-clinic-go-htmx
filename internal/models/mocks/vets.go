package mocks

import "pet-clinic.bonglee.com/internal/models"

type VetModel struct {
}

func (model *VetModel) Insert(firstName, lastName string) (int, error) {

	return 0, nil
}

func (model *VetModel) GetVetsPageLen(pageSize int) (int, error) {

	return 1, nil
}

func (model *VetModel) GetVets(page, pageSize int) ([]models.Vet, error) {

	vets := []models.Vet{}

	return vets, nil
}

func (model *VetModel) GetById(id int) (models.Vet, error) {
	vet := models.Vet{}

	return vet, nil
}

func (model *VetModel) Update(id int, firstName, lastName string) error {

	return nil
}

func (model *VetModel) Remove(id int) error {

	return nil
}
