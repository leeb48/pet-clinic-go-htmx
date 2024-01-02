package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"pet-clinic.bonglee.com/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	data := templateData{}
	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) ownerList(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	data := templateData{}
	app.render(w, r, http.StatusOK, "owner-list.html", data)
}

func (app *application) ownerCreate(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {
	data := templateData{}
	app.render(w, r, http.StatusOK, "owner-create.html", data)
}

type OwnerCreateReq struct {
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	PetName      []string `json:"petName"`
	PetType      []string `json:"petType"`
	PetBirthdate []string `json:"petBirthdate"`
}

func (app *application) ownerCreatePost(w http.ResponseWriter, r *http.Request, pr httprouter.Params) {

	// todo: prototype of how to parse incoming new owner request
	var newOwner OwnerCreateReq
	err := json.NewDecoder(r.Body).Decode(&newOwner)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}

	owner := models.Owner{
		FirstName: newOwner.FirstName,
		LastName:  newOwner.LastName,
	}

	pets := []models.Pet{}

	for i := 0; i < len(newOwner.PetName); i++ {

		birthdate, err := time.Parse("01/02/2006", newOwner.PetBirthdate[i])
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		pets = append(pets, models.Pet{
			Name:      newOwner.PetName[i],
			PetType:   newOwner.PetType[i],
			Birthdate: birthdate,
		})
	}

	fmt.Printf("Owner: %v\nPets: %v", owner, pets)
}
