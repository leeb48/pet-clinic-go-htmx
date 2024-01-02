package models

import "time"

type Pet struct {
	Name      string
	PetType   string
	Birthdate time.Time
}

type PetType struct {
	Name string
}
