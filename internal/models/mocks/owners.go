package mocks

type OwnerModel struct{}

func (model *OwnerModel) Insert(firstName, lastName, addr, state, city, phone, email, birthdate string) (int, error) {

	return 0, nil
}
