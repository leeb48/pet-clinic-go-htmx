package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"pet-clinic.bonglee.com/internal/models"
	"pet-clinic.bonglee.com/internal/models/customErrors"
	"pet-clinic.bonglee.com/internal/validator"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := templateData{}
	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) ownerList(w http.ResponseWriter, r *http.Request) {
	data := templateData{}
	app.render(w, r, http.StatusOK, "owner-list.html", data)
}

type newOwnerForm struct {
	FirstName           string   `json:"firstName"`
	LastName            string   `json:"lastName"`
	Address             string   `json:"address"`
	State               string   `json:"state"`
	City                string   `json:"city"`
	Phone               string   `json:"phone"`
	Email               string   `json:"email"`
	Birthdate           string   `json:"birthdate"`
	PetName             []string `json:"petName"`
	PetType             []string `json:"petType"`
	PetBirthdate        []string `json:"petBirthdate"`
	ValidPetTypes       []string
	validator.Validator `form:"-"`
}

func (app *application) ownerCreate(w http.ResponseWriter, r *http.Request) {

	petTypes, err := app.petTypes.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Form = newOwnerForm{
		ValidPetTypes: petTypes,
	}

	app.render(w, r, http.StatusOK, "owner-create.html", data)
}

// todo refactor to make it cleaner
func (app *application) ownerCreatePost(w http.ResponseWriter, r *http.Request) {
	var form newOwnerForm

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}

	form.Validator.CheckField(validator.NotBlank(form.FirstName), "firstName", "First name cannot be empty")
	form.Validator.CheckField(validator.NotBlank(form.LastName), "lastName", "Last name cannot be empty")
	form.Validator.CheckField(validator.NotBlank(form.Address), "address", "Address cannot be empty")
	form.Validator.CheckField(validator.NotBlank(form.State), "state", "State cannot be empty")
	form.Validator.CheckField(validator.NotBlank(form.City), "city", "City cannot be empty")
	form.Validator.CheckField(validator.NotBlank(form.Phone), "phone", "Phone cannot be empty")
	form.Validator.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be empty")
	form.Validator.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Email invalid")
	form.Validator.CheckField(validator.NotBlank(form.Birthdate), "birthdate", "Birthdate cannot be empty")

	if !form.Validator.Valid() {

		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "owner-create.html", data)
		return
	}

	ownerId, err := app.owners.Insert(form.FirstName, form.LastName, form.Address, form.State, form.City, form.Phone, form.Email, form.Birthdate)
	if err != nil {

		app.logger.Error(err.Error())

		data := app.newTemplateData(r)
		data.Form = form
		data.Alert = Alert{Msg: "Owner creation error", MsgType: DANGER}

		app.render(w, r, http.StatusUnprocessableEntity, "owner-create.html", data)
		return
	}

	pets := []models.Pet{}

	for i := 0; i < len(form.PetName); i++ {

		if form.PetName[i] == "" || form.PetType[i] == "" || form.PetBirthdate[i] == "" {
			continue
		}

		birthdate, err := time.Parse("2006-01-02", (form.PetBirthdate[i]))
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		petTypeId, err := app.petTypes.GetIdFromPetType(form.PetType[i])
		if err != nil {
			app.logger.Error(err.Error())
			petTypeId = app.cfg.DefaultPetType
		}

		pets = append(pets, models.Pet{
			Name:      form.PetName[i],
			Birthdate: birthdate,
			PetTypeId: petTypeId,
			OwnerId:   ownerId,
		})
	}

	for _, pet := range pets {
		err := app.pets.Insert(pet.Name, pet.Birthdate, pet.PetTypeId, pet.OwnerId)
		if err != nil {
			app.logger.Error(err.Error())
		}
	}

	app.session.Put(r.Context(), FLASH_MSG, "User created")

	http.Redirect(w, r, "/owner/create", http.StatusSeeOther)
}

type newPetTypeForm struct {
	NewPetType          string `json:"newPetType"`
	validator.Validator `form:"-"`
}

func (app *application) adminPage(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = newPetTypeForm{}
	app.render(w, r, http.StatusOK, "admin.html", data)
}

func (app *application) newPetTypePost(w http.ResponseWriter, r *http.Request) {
	var form newPetTypeForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.NewPetType = strings.ToUpper(form.NewPetType)

	form.CheckField(validator.NotBlank(form.NewPetType), "newPetType", "Pet type cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "admin.html", data)
		return
	}

	err = app.petTypes.Insert(form.NewPetType)
	if err != nil {

		app.logger.Error(err.Error())

		if errors.Is(err, customErrors.ErrDuplicatePetType) {
			form.AddFieldError("newPetType", "Duplicate pet type")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "admin.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
