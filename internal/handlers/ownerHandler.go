package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pet-clinic.bonglee.com/internal/app"
	"pet-clinic.bonglee.com/internal/models"
	"pet-clinic.bonglee.com/internal/validator"
)

type OwnerHandler struct {
	*app.App
}

func NewOwnerHandler(app *app.App) *OwnerHandler {
	return &OwnerHandler{
		app,
	}
}

func (handler *OwnerHandler) home(w http.ResponseWriter, r *http.Request) {
	data := app.TemplateData{}
	handler.Render(w, r, http.StatusOK, "home.html", data)
}

type ownerListForm struct {
	PageLen int
	Owners  []models.Owner
}

func (handler *OwnerHandler) ownerList(w http.ResponseWriter, r *http.Request) {

	pageSize := r.URL.Query().Get("pageSize")

	if strings.TrimSpace(pageSize) == "" {
		pageSize = "10"
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		handler.ClientError(w, http.StatusBadRequest)
	}

	pageLen, err := handler.Owners.GetOwnersPageLen(pageSizeInt)
	if err != nil {
		handler.ServerError(w, r, err)
	}

	owners, err := handler.Owners.GetOwners(1, pageSizeInt)
	if err != nil {
		handler.ServerError(w, r, err)
	}

	data := app.TemplateData{}
	data.Form = ownerListForm{
		PageLen: pageLen,
		Owners:  owners,
	}
	handler.Render(w, r, http.StatusOK, "owner-list.html", data)
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

func (handler *OwnerHandler) ownerCreate(w http.ResponseWriter, r *http.Request) {

	petTypes, err := handler.PetTypes.GetAll()
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	data := handler.NewTemplateData(r)
	data.Form = newOwnerForm{
		ValidPetTypes: petTypes,
	}

	handler.Render(w, r, http.StatusOK, "owner-create.html", data)
}

func (handler *OwnerHandler) ownerCreatePost(w http.ResponseWriter, r *http.Request) {
	var form newOwnerForm

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		handler.Logger.Error(err.Error())
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

		data := handler.NewTemplateData(r)
		data.Form = form
		handler.Render(w, r, http.StatusUnprocessableEntity, "owner-create.html", data)
		return
	}

	ownerId, err := handler.Owners.Insert(form.FirstName, form.LastName, form.Address, form.State, form.City, form.Phone, form.Email, form.Birthdate)
	if err != nil {

		handler.Logger.Error(err.Error())

		data := handler.NewTemplateData(r)
		data.Form = form
		data.Alert = app.Alert{Msg: "Owner creation error", MsgType: app.DANGER}

		handler.Render(w, r, http.StatusUnprocessableEntity, "owner-create.html", data)
		return
	}

	pets := []models.Pet{}

	for i := 0; i < len(form.PetName); i++ {

		if form.PetName[i] == "" || form.PetType[i] == "" || form.PetBirthdate[i] == "" {
			continue
		}

		birthdate, err := time.Parse("2006-01-02", (form.PetBirthdate[i]))
		if err != nil {
			handler.ServerError(w, r, err)
			return
		}

		petTypeId, err := handler.PetTypes.GetIdFromPetType(form.PetType[i])
		if err != nil {
			handler.Logger.Error(err.Error())
			petTypeId = handler.Cfg.DefaultPetType
		}

		pets = append(pets, models.Pet{
			Name:      form.PetName[i],
			Birthdate: birthdate,
			PetTypeId: petTypeId,
			OwnerId:   ownerId,
		})
	}

	for _, pet := range pets {
		err := handler.Pets.Insert(pet.Name, pet.Birthdate, pet.PetTypeId, pet.OwnerId)
		if err != nil {
			handler.Logger.Error(err.Error())
		}
	}

	handler.Session.Put(r.Context(), app.FLASH_MSG, "User created")

	http.Redirect(w, r, "/owner/create", http.StatusSeeOther)
}
