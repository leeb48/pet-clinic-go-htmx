package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"pet-clinic.bonglee.com/internal/app"
	"pet-clinic.bonglee.com/internal/constants"
	"pet-clinic.bonglee.com/internal/constants/alertConstants"
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

	page := r.URL.Query().Get("page")
	if strings.TrimSpace(page) == "" {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		handler.ClientError(w, http.StatusBadRequest)
	}

	pageLen, err := handler.Owners.GetOwnersPageLen(pageSizeInt)
	if err != nil {
		handler.ServerError(w, r, err)
	}

	owners, err := handler.Owners.GetOwners(pageInt, pageSizeInt)
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
	FirstName           string             `json:"firstName"`
	LastName            string             `json:"lastName"`
	Address             string             `json:"address"`
	State               string             `json:"state"`
	City                string             `json:"city"`
	Phone               string             `json:"phone"`
	Email               string             `json:"email"`
	Birthdate           string             `json:"birthdate"`
	Pets                []models.PetDetail `json:"pets"`
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
		data.Alert = app.Alert{Msg: "Owner creation error", MsgType: alertConstants.DANGER}

		handler.Render(w, r, http.StatusUnprocessableEntity, "owner-create.html", data)
		return
	}

	for _, pet := range form.Pets {

		petTypeId, err := handler.PetTypes.GetIdFromPetType(pet.PetType)
		if err != nil {
			handler.Logger.Error(err.Error())
			petTypeId = handler.Cfg.DefaultPetType
		}

		err = handler.Pets.Insert(pet.Name, pet.Birthdate, petTypeId, ownerId)
		if err != nil {
			handler.Logger.Error(err.Error())
		}
	}

	handler.Session.Put(r.Context(), alertConstants.FLASH_MSG, "User created")

	http.Redirect(w, r, "/owner/create", http.StatusSeeOther)
}

type ownerDetailForm struct {
	Owner models.Owner
	Pets  []models.PetDetail
}

func (handler *OwnerHandler) ownerDetail(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		handler.ServerError(w, r, err)
	}

	owner, err := handler.Owners.GetOwnerById(id)
	if err != nil {
		handler.ServerError(w, r, err)
	}

	pets, err := handler.Pets.GetPetsByOwnerId(owner.Id)
	if err != nil {
		handler.ServerError(w, r, err)
	}

	data := handler.NewTemplateData(r)
	data.Form = ownerDetailForm{
		Owner: owner,
		Pets:  pets,
	}

	handler.Render(w, r, http.StatusOK, "owner-detail.html", data)
}

type editOwnerForm struct {
	Id                  int
	FirstName           string             `json:"firstName"`
	LastName            string             `json:"lastName"`
	Address             string             `json:"address"`
	State               string             `json:"state"`
	City                string             `json:"city"`
	Phone               string             `json:"phone"`
	Email               string             `json:"email"`
	Birthdate           string             `json:"birthdate"`
	Pets                []models.PetDetail `json:"pets"`
	DeletePetIds        []int              `json:"deletePetIds"`
	ValidPetTypes       []string
	validator.Validator `form:"-"`
}

func (handler *OwnerHandler) ownerEdit(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		handler.ServerError(w, r, err)
	}

	owner, err := handler.Owners.GetOwnerById(id)
	if err != nil {
		handler.ServerError(w, r, err)
	}

	validPetTypes, err := handler.PetTypes.GetAll()
	if err != nil {
		handler.ServerError(w, r, err)
	}

	pets, err := handler.Pets.GetPetsByOwnerId(id)
	if err != nil {
		handler.ServerError(w, r, err)
	}

	data := handler.NewTemplateData(r)
	data.Form = editOwnerForm{
		Id:            id,
		FirstName:     owner.FirstName,
		LastName:      owner.LastName,
		Email:         owner.Email,
		Phone:         owner.Phone,
		Birthdate:     owner.Birthdate.Format(constants.YYYY_MM_DD),
		Address:       owner.Address,
		City:          owner.City,
		State:         owner.State,
		Pets:          pets,
		ValidPetTypes: validPetTypes,
	}

	handler.Render(w, r, http.StatusOK, "owner-edit.html", data)
}

func (handler *OwnerHandler) ownerEditPut(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		handler.ClientError(w, http.StatusBadRequest)
		return
	}

	var form editOwnerForm

	err = json.NewDecoder(r.Body).Decode(&form)
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

		validPetTypes, err := handler.PetTypes.GetAll()
		if err != nil {
			handler.ServerError(w, r, err)
		}
		form.ValidPetTypes = validPetTypes

		data.Form = form
		handler.Render(w, r, http.StatusUnprocessableEntity, "owner-edit.html", data)
		return
	}

	err = handler.Owners.UpdateOwner(id, form.FirstName, form.LastName, form.Address, form.State, form.City, form.Phone, form.Email, form.Birthdate)
	if err != nil {
		handler.ServerError(w, r, err)
		return
	}

	for _, petId := range form.DeletePetIds {
		err := handler.Pets.Remove(petId)
		if err != nil {
			handler.Logger.Error(err.Error())
		}
	}

	for _, pet := range form.Pets {

		petTypeId, err := handler.PetTypes.GetIdFromPetType(pet.PetType)
		if err != nil {
			handler.Logger.Error(err.Error())
			petTypeId = handler.Cfg.DefaultPetType
		}

		if pet.Id == 0 {
			err = handler.Pets.Insert(pet.Name, pet.Birthdate, petTypeId, id)
		} else {
			err = handler.Pets.Update(pet.Id, pet.Name, pet.Birthdate, petTypeId)
		}

		if err != nil {
			handler.Logger.Error(err.Error())
		}
	}

	handler.Session.Put(r.Context(), alertConstants.FLASH_MSG, "Edit Successful")

	w.Header().Add("HX-Redirect", fmt.Sprintf("/owner/detail/%v", id))
}
