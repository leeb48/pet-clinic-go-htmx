package models

import "time"

type PetModel struct {
	Id        int
	Name      string
	PetType   string
	Birthdate time.Time
}
